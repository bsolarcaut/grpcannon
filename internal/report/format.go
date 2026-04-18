package report

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// Format controls how a RunReport is rendered.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Print writes the report to w in the requested format.
func Print(w io.Writer, r RunReport, f Format) error {
	switch f {
	case FormatJSON:
		return printJSON(w, r)
	default:
		return printText(w, r)
	}
}

func printText(w io.Writer, r RunReport) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Target:\t%s\n", r.Target)
	fmt.Fprintf(tw, "Method:\t%s\n", r.Method)
	fmt.Fprintf(tw, "Workers:\t%d\n", r.Workers)
	fmt.Fprintf(tw, "Duration:\t%s\n", r.Duration.Round(1e6))
	fmt.Fprintf(tw, "Total:\t%d\n", r.Summary.Total)
	fmt.Fprintf(tw, "Success:\t%d\n", r.Summary.Success)
	fmt.Fprintf(tw, "Errors:\t%d\n", r.Summary.Errors)
	fmt.Fprintf(tw, "RPS:\t%.2f\n", r.RPS)
	fmt.Fprintf(tw, "P50:\t%s\n", r.Summary.P50)
	fmt.Fprintf(tw, "P95:\t%s\n", r.Summary.P95)
	fmt.Fprintf(tw, "P99:\t%s\n", r.Summary.P99)
	return tw.Flush()
}

func printJSON(w io.Writer, r RunReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
