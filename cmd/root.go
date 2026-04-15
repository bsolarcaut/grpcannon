package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/grpcannon/internal/config"
	"github.com/spf13/cobra"
)

var cfg = config.DefaultConfig()

var rootCmd = &cobra.Command{
	Use:   "grpcannon",
	Short: "A lightweight load-testing CLI for gRPC services",
	Long: `grpcannon fires configurable concurrent gRPC requests at a target
and reports latency histograms, throughput, and error rates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}
		fmt.Printf("Target:       %s\n", cfg.Target)
		fmt.Printf("Method:       %s\n", cfg.Method)
		fmt.Printf("Concurrency:  %d\n", cfg.Concurrency)
		if cfg.Duration > 0 {
			fmt.Printf("Duration:     %s\n", cfg.Duration)
		} else {
			fmt.Printf("Requests:     %d\n", cfg.TotalRequests)
		}
		fmt.Println("(runner not yet implemented)")
		return nil
	},
}

func init() {
	f := rootCmd.Flags()
	f.StringVarP(&cfg.Target, "target", "t", "", "gRPC server address (host:port) [required]")
	f.StringVarP(&cfg.Method, "method", "m", "", "fully-qualified gRPC method (pkg.Svc/Method) [required]")
	f.IntVarP(&cfg.Concurrency, "concurrency", "c", cfg.Concurrency, "number of concurrent workers")
	f.IntVarP(&cfg.TotalRequests, "requests", "n", cfg.TotalRequests, "total requests to send (ignored when --duration is set)")
	f.DurationVarP(&cfg.Duration, "duration", "d", 0, "test duration (e.g. 30s); overrides --requests")
	f.DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "per-request timeout")
	f.BoolVar(&cfg.Insecure, "insecure", cfg.Insecure, "disable TLS certificate verification")
	f.StringToStringVar(&cfg.Metadata, "metadata", cfg.Metadata, "gRPC metadata key=value pairs")
	f.StringVarP(&cfg.PayloadJSON, "payload", "p", "", "JSON-encoded request payload")

	_ = rootCmd.MarkFlagRequired("target")
	_ = rootCmd.MarkFlagRequired("method")
}

// Execute is the entry-point called by main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// keep time import used via cfg default
var _ = time.Second
