package cleanUtils

import (
	"cheeseshadow/libenforcer/types"
	"fmt"
	"os"
	"path/filepath"
)

func CleanLibrary(libPath string, tracks []types.TrackTransform) (err error) {
	whitelist := make(map[string]bool)
	for _, track := range tracks {
		whitelist[filepath.Join(libPath, track.TrackPath())] = true
	}

	if err = cleanFiles(libPath, whitelist); err != nil {
		return
	}

	if err = cleanFolders(libPath); err != nil {
		return
	}

	return
}

func cleanFiles(root string, whitelist map[string]bool) (err error) {
	files, err := os.ReadDir(root)
	if err != nil {
		return
	}

	for _, file := range files {
		filePath := filepath.Join(root, file.Name())
		if file.IsDir() {
			if err = cleanFiles(filePath, whitelist); err != nil {
				return
			}
		} else {
			if !whitelist[filePath] {
				fmt.Println("Removing file:", filePath)
				if err = os.Remove(filePath); err != nil {
					return
				}
			}
		}
	}

	return
}

func cleanFolders(libPath string) (err error) {
	files, err := os.ReadDir(libPath)
	if err != nil {
		return
	}

	if len(files) == 0 {
		fmt.Println("Removing empty folder:", libPath)
		err = os.Remove(libPath)

		return
	}

	for _, file := range files {
		filePath := filepath.Join(libPath, file.Name())
		if file.IsDir() {
			if err = cleanFolders(filePath); err != nil {
				return
			}
		}

	}

	return
}
