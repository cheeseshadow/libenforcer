package main

import (
	"cheeseshadow/libenforcer/cleanUtils"
	"cheeseshadow/libenforcer/trackUtils"
	"cheeseshadow/libenforcer/types"
	"cheeseshadow/libenforcer/utils"
	"errors"
	"fmt"
	"github.com/schollz/progressbar"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)
import "flag"

func main() {
	libpath := flag.String("f", "", "path to library")
	threads := flag.Int("t", 1, "number of threads to use")
	flag.Parse()

	if *libpath == "" {
		fmt.Println("Please specify a library path using -f")
		return
	}

	tracks, errs := buildTargetChange(*libpath, *threads)
	if len(errs.Get()) > 0 {
		fmt.Println("Errors:")
		for _, err := range errs.Get() {
			fmt.Println(err)
		}
	}

	enforceChange(*libpath, tracks.Get())

	fmt.Println("\nI'll make a slight pause here...")
	time.Sleep(1 * time.Second)

	err := cleanUtils.CleanLibrary(*libpath, tracks.Get())
	if err != nil {
		fmt.Println(err)
	}
}

func buildTargetChange(libPath string, maxGoroutinesAtOnce int) (tracks types.ConcurrentCollection[types.TrackTransform], errs types.ConcurrentCollection[error]) {
	fmt.Println("Building target change...")

	files, err := os.ReadDir(libPath)
	if err != nil {
		fmt.Println("Failed to find the libpath:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		fmt.Println("Processing", file.Name())

		filePath := filepath.Join(libPath, file.Name())
		dirTracks, retryTracks, asyncErrs := traverseConcurrently(filePath, maxGoroutinesAtOnce)
		dirErrs := asyncErrs.Get()

		if len(retryTracks.Get()) > 0 {
			fmt.Printf("Retrying %v failed tracks...\n", len(retryTracks.Get()))
			syncTracks, syncErrs := handleTracksSynchronously(retryTracks.Get())

			dirTracks.AppendAll(syncTracks)
			dirErrs = append(dirErrs, syncErrs...)
		}
		fmt.Println("Errors:", len(dirErrs))

		enforceChange(libPath, tracks.Get())
		fmt.Println("Change enforced...")

		tracks.AppendAll(dirTracks.Get())
		errs.AppendAll(dirErrs)
	}

	return
}

func handleTracksSynchronously(trackPaths []string) (tracks []types.TrackTransform, errs []error) {
	for _, trackPath := range trackPaths {
		albumPath, trackName, err := trackUtils.HandleTrack(trackPath)
		if err != nil {
			errs = append(errs, err)
		} else {
			tracks = append(tracks, types.TrackTransform{
				OriginalPath: trackPath,
				AlbumPath:    albumPath,
				TrackName:    trackName,
			})
		}
	}

	return
}

func traverseConcurrently(
	libPath string,
	maxGoroutinesAtOnce int,
) (tracks types.ConcurrentCollection[types.TrackTransform], retryTracks types.ConcurrentCollection[string], errs types.ConcurrentCollection[error]) {
	var wg sync.WaitGroup
	semaphore := make(chan bool, maxGoroutinesAtOnce)

	wg.Add(1)
	go traverse(libPath, &tracks, &errs, &retryTracks, &wg, semaphore)
	wg.Wait()

	return
}

func enforceChange(libPath string, tracks []types.TrackTransform) {
	trackCount := len(tracks)
	bar := progressbar.New(trackCount)

	for _, track := range tracks {
		bar.Add(1)

		fullAlbumPath := filepath.Join(libPath, track.AlbumPath)
		fullTrackPath := filepath.Join(libPath, track.TrackPath())
		if track.OriginalPath == fullTrackPath {
			continue
		}

		if err := os.MkdirAll(fullAlbumPath, 0755); err != nil {
			fmt.Println(err)
		}

		if err := os.Rename(track.OriginalPath, fullTrackPath); err != nil {
			fmt.Println(err)
		}
	}
}

func traverse(
	path string,
	tracks *types.ConcurrentCollection[types.TrackTransform],
	errs *types.ConcurrentCollection[error],
	retryTracks *types.ConcurrentCollection[string],
	wg *sync.WaitGroup,
	semaphore chan bool,
) {
	semaphore <- true
	defer func() {
		wg.Done()
		<-semaphore
	}()

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Failed to find the libpath:", err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		if file.IsDir() {
			wg.Add(1)
			go traverse(filepath.Join(filePath), tracks, errs, retryTracks, wg, semaphore)
		} else {
			fileType, err := utils.GetFileContentType(filePath)
			if err != nil {
				errs.Add(errors.New(fmt.Sprintf("Failed to handle track %s: %s", filePath, err)))
			}
			if fileType != "audio/mpeg" {
				continue
			}

			albumPath, trackName, err := trackUtils.HandleTrack(filePath)
			if err != nil {
				if strings.Contains(err.Error(), "device not configured") ||
					strings.Contains(err.Error(), "bad file descriptor") ||
					strings.Contains(err.Error(), "input/output error") {
					retryTracks.Add(filePath)
				} else {
					errs.Add(errors.New(fmt.Sprintf("Failed to handle track %s: %s", filePath, err)))
				}
			} else {
				tracks.Append(types.TrackTransform{
					OriginalPath: filePath,
					AlbumPath:    albumPath,
					TrackName:    trackName,
				})
			}
		}
	}
}
