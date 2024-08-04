/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	vegeta "github.com/tsenart/vegeta/lib"
)

func NewPostCustomTargeter() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "POST"
		tgt.URL = "http://server-savings-service-tpe.service.i.gojek.gcp/internal/debug/benchmark"

		return nil
	}
}

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Post Benchmark",
	Long:  `Post benchmark.`,
	Run: func(cmd *cobra.Command, args []string) {
		rate := vegeta.Rate{Freq: 100, Per: time.Second}
		duration := 60 * time.Second
		targeter := NewPostCustomTargeter()
		attacker := vegeta.NewAttacker()

		var metrics vegeta.Metrics
		for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
			metrics.Add(res)
		}
		metrics.Close()
		fmt.Printf("Rate: %vreq/s\n", metrics.Rate)
		fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
		fmt.Printf("95th percentile: %s\n", metrics.Latencies.P95)
		fmt.Printf("50 percentile: %s\n", metrics.Latencies.P50)
		fmt.Printf("Max: %s\n", metrics.Latencies.Max)
		fmt.Printf("Mean: %s\n", metrics.Latencies.Mean)
		fmt.Println("Status Error", metrics.Errors)
		fmt.Println("Status Success", metrics.StatusCodes)
	},
}

func init() {
	rootCmd.AddCommand(postCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// postCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// postCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
