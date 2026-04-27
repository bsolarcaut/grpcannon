// Package shedder provides adaptive load shedding for grpcannon.
//
// A Shedder tracks a sliding window of call outcomes and begins rejecting
// new requests (returning ErrShed from Allow) when the error rate in that
// window exceeds a configurable threshold.  After a cooldown period the
// shedder re-evaluates and may resume allowing traffic.
//
// Typical usage:
//
//	s := shedder.New(shedder.Default())
//	if err := s.Allow(); err != nil {
//		// drop / return 503
//	}
//	err := doRPC()
//	s.Record(err != nil)
package shedder
