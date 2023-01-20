package main

import (
	"cheeseshadow/libenforcer/cleanUtils"
	"cheeseshadow/libenforcer/trackUtils"
	"cheeseshadow/libenforcer/types"
	"fmt"
	"os"
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
	enforceChange(*libpath, tracks)
	cleanUtils.CleanLibrary(*libpath, tracks)

	fmt.Println(tracks)
	fmt.Println(errs)
}

func buildTargetChange(libPath string) (tracks []types.TrackTransform, errs []error) {
	traverse(libPath, &tracks, &errs)
	return
}

func enforceChange(libPath string, tracks []types.TrackTransform) {
	for _, track := range tracks {
		fullAlbumPath := libPath + "/" + track.AlbumPath
		fullTrackPath := fullAlbumPath + "/" + track.TrackName
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

func traverse(path string, tracks *[]types.TrackTransform, errs *[]error) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Failed to find the libpath:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			traverse(path+"/"+file.Name(), tracks, errs)
		} else {
			albumPath, trackName, err := trackUtils.HandleTrack(path + "/" + file.Name())
			if err != nil {
				*errs = append(*errs, err)
			} else {
				*tracks = append(*tracks, types.TrackTransform{
					OriginalPath: path + "/" + file.Name(),
					AlbumPath:    albumPath,
					TrackName:    trackName,
				})
			}
		}
	}
}
