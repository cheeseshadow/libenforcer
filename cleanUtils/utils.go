package cleanUtils

import (
	"cheeseshadow/libenforcer/types"
	"os"
)

func CleanLibrary(libPath string, tracks []types.TrackTransform) (err error) {
	whitelist := make(map[string]bool)
	for _, track := range tracks {
		whitelist[libPath+"/"+track.AlbumPath+"/"+track.TrackName] = true
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
		if file.IsDir() {
			if err = cleanFiles(root+"/"+file.Name(), whitelist); err != nil {
				return
			}
		} else {
			filePath := root + "/" + file.Name()
			if !whitelist[filePath] {
				if err = os.Remove(root + "/" + file.Name()); err != nil {
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
		err = os.Remove(libPath)

		return
	}

	for _, file := range files {
		if file.IsDir() {
			if err = cleanFolders(libPath + "/" + file.Name()); err != nil {
				return
			}
		}

	}

	return
}
