// Package window implements a sliding-window counter used to track
// request throughput over a configurable rolling time interval.
//
// Usage:
//
//	w := window.New(10, time.Second) // 10-second window
//	w.Add(1)
//	fmt.Println(w.Rate()) // requests/sec
package window
