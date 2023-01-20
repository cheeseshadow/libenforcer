package types

import (
	"path/filepath"
)

type TrackTransform struct {
	OriginalPath string
	AlbumPath    string
	TrackName    string
}

func (t *TrackTransform) TrackPath() string {
	return filepath.Join(t.AlbumPath, t.TrackName)
}

func (t *TrackTransform) String() string {
	return t.OriginalPath + " -> " + t.TrackPath()
}
