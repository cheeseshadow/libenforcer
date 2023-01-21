package cleanUtils

import (
	"cheeseshadow/libenforcer/types"
	"cheeseshadow/libenforcer/utils"
	"fmt"
	"os"
	"path/filepath"
)

type deletedCounter struct {
	Value int
}

func CleanLibrary(libPath string, tracks []types.TrackTransform) (err error) {
	fmt.Println("Cleaning library...")

	whitelist := make(map[string]bool)
	for _, track := range tracks {
		whitelist[filepath.Join(libPath, track.TrackPath())] = true
	}

	for {
		var deletedFiles deletedCounter
		fileErrs := cleanFiles(libPath, whitelist, &deletedFiles)
		if len(fileErrs) > 0 {
			fmt.Println("Some files could not be deleted:")
			for _, fileErr := range fileErrs {
				fmt.Println(fileErr)
			}
		}

		var deletedDirs deletedCounter
		if err = cleanFolders(libPath, &deletedDirs); err != nil {
			return
		}

		if deletedFiles.Value == 0 && deletedDirs.Value == 0 {
			break
		}
	}

	return
}

func cleanFiles(root string, whitelist map[string]bool, counter *deletedCounter) (errs []error) {
	files, err := os.ReadDir(root)
	if err != nil {
		errs = append(errs, err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(root, file.Name())
		if !file.IsDir() {
			if _, ok := whitelist[filePath]; !ok {
				if utils.CheckIfFileExists(filePath) {
					fmt.Println("Removing file:", filePath)
					counter.Value++
					if err := os.Remove(filePath); err != nil {
						errs = append(errs, err)
					}
				}
			}
		}
	}

	for _, file := range files {
		filePath := filepath.Join(root, file.Name())
		if file.IsDir() {
			errs = append(errs, cleanFiles(filePath, whitelist, counter)...)
		}
	}

	return
}

func cleanFolders(libPath string, counter *deletedCounter) (err error) {
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
