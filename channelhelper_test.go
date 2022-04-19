package channelhelpers

import (
	"fmt"
	"testing"
	"time"
)

func TestBatchEvents(t *testing.T) {
	const numEvents = 1000
	ch := make(chan string, numEvents)
	for i := 0; i < numEvents; i++ {
		ch <- fmt.Sprintf("%d", i+1)
	}
	close(ch)
	counter := 0
	sum := 0
	for e := range BatchEvents(ch, 100, 2*time.Second) {
		length := len(e)
		sum += length
		if length != 100 {
			t.Errorf("length %d != 100", length)
		}
		counter += 1
	}
	if counter != 10 {
		t.Errorf("counter %d != 10", counter)
	}

	fmt.Printf("sum: %d\n", sum)

	if sum != numEvents {
		t.Errorf("sum %d != numEvents %d", sum, numEvents)
	}
}
