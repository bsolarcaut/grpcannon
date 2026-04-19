package fanout_test

import (
	"sync"
	"testing"

	"github.com/example/grpcannon/internal/fanout"
)

func TestNew_DefaultsBuffer(t *testing.T) {
	fo := fanout.New[int](0)
	if fo == nil {
		t.Fatal("expected non-nil fanout")
	}
}

func TestSubscribe_ReceivesAllValues(t *testing.T) {
	fo := fanout.New[int](8)
	s1 := fo.Subscribe()
	s2 := fo.Subscribe()

	src := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		src <- i
	}
	close(src)

	fo.Run(src)

	for _, sink := range []fanout.Sink[int]{s1, s2} {
		sum := 0
		for v := range sink {
			sum += v
		}
		if sum != 15 {
			t.Errorf("expected sum 15, got %d", sum)
		}
	}
}

func TestSubscribe_ChannelsClosedAfterRun(t *testing.T) {
	fo := fanout.New[string](4)
	sink := fo.Subscribe()

	src := make(chan string)
	close(src)
	fo.Run(src)

	_, open := <-sink
	if open {
		t.Error("expected sink to be closed")
	}
}

func TestFanout_ConcurrentSubscribers(t *testing.T) {
	const n = 10
	fo := fanout.New[int](16)

	sinks := make([]fanout.Sink[int], n)
	for i := range sinks {
		sinks[i] = fo.Subscribe()
	}

	src := make(chan int, n)
	for i := 0; i < n; i++ {
		src <- i
	}
	close(src)
	fo.Run(src)

	var wg sync.WaitGroup
	for _, s := range sinks {
		wg.Add(1)
		go func(sink fanout.Sink[int]) {
			defer wg.Done()
			count := 0
			for range sink {
				count++
			}
			if count != n {
				t.Errorf("expected %d items, got %d", n, count)
			}
		}(s)
	}
	wg.Wait()
}
