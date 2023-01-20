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

type ConcurrentWaitGroup struct {
	mutex sync.Mutex
	wg    sync.WaitGroup
}

func (c *ConcurrentWaitGroup) Add(delta int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.wg.Add(delta)
}

func (c *ConcurrentWaitGroup) Done() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.wg.Done()
}

func (c *ConcurrentWaitGroup) Wait() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.wg.Wait()
}
