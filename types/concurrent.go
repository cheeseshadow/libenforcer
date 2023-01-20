package types

import "sync"

type ConcurrentTrackCollection struct {
	mutex  sync.Mutex
	tracks []TrackTransform
}

func (c *ConcurrentTrackCollection) Add(track TrackTransform) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.tracks = append(c.tracks, track)
}

func (c *ConcurrentTrackCollection) Get() []TrackTransform {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.tracks
}

func (c *ConcurrentTrackCollection) Append(track TrackTransform) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.tracks = append(c.tracks, track)
}
