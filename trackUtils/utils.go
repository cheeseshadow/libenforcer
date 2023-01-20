package trackUtils

import (
	"fmt"
	"github.com/bogem/id3v2"
	"path/filepath"
	"strconv"
	"strings"
)

func HandleTrack(trackPath string) (albumPath string, trackName string, err error) {
	tag, err := id3v2.Open(trackPath, id3v2.Options{Parse: true})
	if err != nil {
		return
	}

	trackNumber, _, err := parseSet(tag.GetTextFrame(tag.CommonID("Track number/Position in set")).Text)
	if err != nil {
		return
	}

	discNumber, discCount, err := parseSet(tag.GetTextFrame(tag.CommonID("Part of a set")).Text)
	if err != nil {
		return
	}

	extension := filepath.Ext(trackPath)
	albumPath, trackName = buildPath(tag.Artist(), tag.Album(), tag.Title(), tag.Year(), trackNumber, discNumber, discCount, extension)

	return
}

func buildPath(artist string, album string, title string, year string, trackNumber int, discNumber int, discCount int, extension string) (string, string) {
	trackName := fmt.Sprintf("%02d - %s%s", trackNumber, title, extension)
	if discCount > 1 {
		trackName = strconv.Itoa(discNumber) + trackName
	}

	albumPath := fmt.Sprintf("%s/%s - %s", artist, year, album)

	return albumPath, trackName
}

func parseSet(set string) (number int, count int, err error) {
	setData := strings.Split(set, "/")

	number, err = strconv.Atoi(setData[0])
	if err != nil {
		return
	}

	count, err = strconv.Atoi(setData[1])
	if err != nil {
		return
	}

	return
}
