package main

import (
	"cheeseshadow/libenforcer/cleanUtils"
	"cheeseshadow/libenforcer/trackUtils"
	"cheeseshadow/libenforcer/types"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)
import "flag"

func main() {
	libpath := flag.String("f", "", "path to library")
	flag.Parse()

	if *libpath == "" {
		fmt.Println("Please specify a library path using -f")
		return
	}

	tracks, errs := buildTargetChange(*libpath)
	if len(errs.Get()) > 0 {
		fmt.Println("Errors:")
		for _, err := range errs.Get() {
			fmt.Println(err)
		}
	}

	enforceChange(*libpath, tracks.Get())
	err := cleanUtils.CleanLibrary(*libpath, tracks.Get())
	if err != nil {
		fmt.Println(err)
	}
}

func buildTargetChange(libPath string) (tracks types.ConcurrentCollection[types.TrackTransform], errs types.ConcurrentCollection[error]) {
	fmt.Println("Building target change...")
	maxGoroutinesAtOnce := 5

	var wg sync.WaitGroup
	semaphore := make(chan bool, maxGoroutinesAtOnce)
	wg.Add(1)
	go traverse(libPath, &tracks, &errs, &wg, semaphore)
	wg.Wait()

	return
}

func enforceChange(libPath string, tracks []types.TrackTransform) {
	trackCount := len(tracks)
	for trackNum, track := range tracks {
		fmt.Printf("Handling track %d/%d: %s\n", trackNum+1, trackCount, track.TrackName)

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

func traverse(path string, tracks *types.ConcurrentCollection[types.TrackTransform], errs *types.ConcurrentCollection[error], wg *sync.WaitGroup, semaphore chan bool) {
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
			go traverse(filepath.Join(filePath), tracks, errs, wg, semaphore)
		} else {
			albumPath, trackName, err := trackUtils.HandleTrack(filePath)
			if err != nil {
				errs.Add(errors.New(fmt.Sprintf("Failed to handle track %s: %s", filePath, err)))
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
