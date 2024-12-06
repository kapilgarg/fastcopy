package progress

import (
	"fmt"
	"sync"
	"time"
)

type Tracker struct {
	TotalCopied int64
	FileSize    int64
	mutex       sync.Mutex
	stopTicker  chan bool
}

func NewTracker(fileSize int64) *Tracker {
	return &Tracker{
		FileSize:   fileSize,
		stopTicker: make(chan bool),
	}
}

func (t *Tracker) AddCopied(bytes int64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.TotalCopied += bytes
}

func StartProgressTicker(t *Tracker) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				printProgress(t)
			case <-t.stopTicker:
				printProgress(t) // Ensure final update
				return
			}
		}
	}()
}

func StopProgressTicker(t *Tracker) {
	close(t.stopTicker)
}

func printProgress(t *Tracker) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	percent := int(float64(t.TotalCopied) / float64(t.FileSize) * 100)
	fmt.Printf("\rCopying... %d%%", percent)
}
