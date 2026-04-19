// Package fanout provides a generic broadcast primitive that fans a single
// producer channel out to an arbitrary number of subscriber channels.
//
// Typical usage:
//
//	fo := fanout.New[MyResult](32)
//	sink1 := fo.Subscribe()
//	sink2 := fo.Subscribe()
//	go fo.Run(resultsCh)
//
// Each subscriber receives every value emitted by the source channel.
// All subscriber channels are closed when the source channel closes.
package fanout
