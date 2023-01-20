package cleanUtils

import (
	"cheeseshadow/libenforcer/types"
	"fmt"
	"os"
	"path/filepath"
)

type deletedCounter struct {
	Value int
}

func CleanLibrary(libPath string, tracks []types.TrackTransform) (err error) {
	whitelist := make(map[string]bool)
	for _, track := range tracks {
		whitelist[filepath.Join(libPath, track.TrackPath())] = true
	}

	for {
		var deletedFiles deletedCounter
		if err = cleanFiles(libPath, whitelist, deletedFiles); err != nil {
			return
		}

		var deletedDirs deletedCounter
		if err = cleanFolders(libPath, deletedDirs); err != nil {
			return
		}

		if deletedFiles.Value == 0 && deletedDirs.Value == 0 {
			break
		}
	}

	return
}

func cleanFiles(root string, whitelist map[string]bool, counter deletedCounter) (err error) {
	files, err := os.ReadDir(root)
	if err != nil {
		return
	}

	for _, file := range files {
		filePath := filepath.Join(root, file.Name())
		if file.IsDir() {
			if err = cleanFiles(filePath, whitelist, counter); err != nil {
				return
			}
		} else {
			if !whitelist[filePath] {
				fmt.Println("Removing file:", filePath)
				counter.Value++
				if err = os.Remove(filePath); err != nil {
					return
				}
			}
		}
	}

	return
}

func cleanFolders(libPath string, counter deletedCounter) (err error) {
	files, err := os.ReadDir(libPath)
	if err != nil {
		return
	}

	if len(files) == 0 {
		fmt.Println("Removing empty folder:", libPath)
		counter.Value++
		err = os.Remove(libPath)

		return
	}

	for _, file := range files {
		filePath := filepath.Join(libPath, file.Name())
		if file.IsDir() {
			if err = cleanFolders(filePath, counter); err != nil {
				return
			}
		}

	}

	return
}
