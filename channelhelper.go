package channelhelpers

import "time"

func BatchEvents[T any](values chan T, maxItems int, maxTimeout time.Duration) <-chan []T {
	batchedChan := make(chan []T)

	dupSlice := func(src []T) []T {
		copyBatch := make([]T, len(src))
		copy(copyBatch, src)
		return copyBatch
	}

	go func() {
		defer close(batchedChan)
		ticker := time.NewTicker(maxTimeout)
		defer ticker.Stop()
		batch := make([]T, 0, maxItems)
		for {
			select {
			case e, more := <-values:
				if !more {
					if len(batch) > 0 {
						batchedChan <- dupSlice(batch)
					}
					return
				} else if len(batch) < maxItems {
					batch = append(batch, e)
				} else {
					batchedChan <- dupSlice(batch)
					batch = make([]T, 0, maxItems)
					batch = append(batch, e)
				}
			case <-ticker.C:
				if len(batch) > 0 {
					batchedChan <- dupSlice(batch)
					batch = make([]T, 0, maxItems)
				}
			}
		}
	}()
	return batchedChan
}
