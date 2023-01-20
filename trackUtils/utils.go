package trackUtils

import (
	"fmt"
	"github.com/bogem/id3v2"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func HandleTrack(trackPath string) (albumPath string, trackName string, err error) {
	tag, err := openWithRetry(trackPath)
	if err != nil {
		return
	}
	defer tag.Close()

	extension := filepath.Ext(trackPath)

	trackNumber, _, err := parseSet(tag.GetTextFrame(tag.CommonID("Track number/Position in set")).Text)
	if err != nil {
		return
	}

	discNumber, discCount, err := parseSet(tag.GetTextFrame(tag.CommonID("Part of a set")).Text)
	if err != nil {
		return
	}

	albumPath, trackName = buildPath(
		sanitizeArtist(tag.Artist()),
		sanitizeName(tag.Album()),
		sanitizeName(tag.Title()),
		tag.Year(),
		trackNumber,
		discNumber,
		discCount,
		extension,
	)

	return
}

func openWithRetry(path string) (tag *id3v2.Tag, err error) {
	maxRetryCount := 5
	for i := 0; i < maxRetryCount; i++ {
		tag, err = id3v2.Open(path, id3v2.Options{Parse: true})
		if err == nil {
			return tag, nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return
}

func buildPath(artist string, album string, title string, year string, trackNumber int, discNumber int, discCount int, extension string) (string, string) {
	trackName := fmt.Sprintf("%02d - %s%s", trackNumber, title, extension)
	if discCount > 1 {
		trackName = strconv.Itoa(discNumber) + trackName
	}

	albumPath := fmt.Sprintf("%s/%s - %s", artist, year, album)

	return albumPath, sanitizeName(trackName)
}

func parseSet(set string) (number int, count int, err error) {
	setData := strings.Split(set, "/")

	if len(setData) == 0 {
		return 0, 0, nil
	}

	number, err = strconv.Atoi(setData[0])
	if err != nil {
		return
	}

	if len(setData) == 1 {
		return number, 1, nil
	}

	count, err = strconv.Atoi(setData[1])
	if err != nil {
		return
	}

	return
}

func sanitizeName(name string) string {
	sanitized := strings.Replace(name, "/", "-", -1)
	sanitized = strings.Replace(sanitized, ":", " -", -1)
	sanitized = strings.Replace(sanitized, "\"", "-", -1)
	return strings.Trim(sanitized, " ")
}

func sanitizeArtist(artist string) string {
	sanitizedArtst := artist

	for _, featTerm := range []string{"feat.", "ft.", "feat", " ft ", "featuring"} {
		if strings.Contains(artist, featTerm) {
			sanitizedArtst = strings.Split(artist, featTerm)[0]
			break
		}
	}

	return sanitizeName(sanitizedArtst)
}
