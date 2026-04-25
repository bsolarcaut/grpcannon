// Package debounce provides a Debouncer that coalesces rapid successive
// calls into a single delayed invocation.
//
// # Usage
//
//	d := debounce.New(200*time.Millisecond, func() {
//		fmt.Println("flushed")
//	})
//	for _, event := range events {
//		d.Trigger() // only the last trigger fires the callback
//	}
//
// Call Flush to fire immediately, or Stop to cancel without firing.
package debounce
