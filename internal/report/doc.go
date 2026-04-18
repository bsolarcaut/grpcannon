// Package report provides a Builder that accumulates per-call results during
// a load-test run and produces a RunReport once the run completes.
//
// Typical usage:
//
//	b := report.New(cfg.Target, cfg.Method, cfg.Concurrency)
//	for _, res := range results {
//		b.Add(res)
//	}
//	r := b.Build()
//	report.Print(os.Stdout, r, report.FormatText)
//
// The RunReport embeds a stats.Summary so all latency percentiles and error
// counts are available for downstream rendering.
package report
