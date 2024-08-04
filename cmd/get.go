/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	vegeta "github.com/tsenart/vegeta/lib"
	"golang.org/x/exp/rand"
)

var customerIDs = []int{
	3612181,
	6132941,
	3331174,
	6055614,
	6241343,
	2334679,
	1254653,
	695385,
	6233386,
	6011242,
	3252604,
	1342151,
	6233819,
	748460,
	4755767,
	221470,
	2557187,
	82910,
	3253273,
	6269775,
	6241264,
	3331184,
	6241057,
	3331155,
	5473266,
	3854843,
	2269345,
	6211728,
	6269839,
	1600504,
	6164610,
	6244592,
	3487084,
	3331159,
	6203468,
	6259014,
	5969130,
	4405729,
	3995611,
}

func getRandomCustomerID(customerIDs []int) int {
	randomIndex := rand.Intn(len(customerIDs))
	return customerIDs[randomIndex]
}
func NewGetCustomTargeter() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "GET"
		customerId := getRandomCustomerID(customerIDs)
		tgt.URL = fmt.Sprintf("http://server-savings-service-tpe.service.i.gojek.gcp/internal/v1/savings/plan_detail/%d", customerId)

		return nil
	}
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get benchmark",
	Long:  `Get benchmark`,
	Run: func(cmd *cobra.Command, args []string) {
		rate := vegeta.Rate{Freq: 100, Per: time.Second}
		duration := 10 * time.Second
		targeter := NewGetCustomTargeter()
		attacker := vegeta.NewAttacker()

		var metrics vegeta.Metrics
		for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
			metrics.Add(res)
		}
		metrics.Close()
		fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
		fmt.Println("Status Error", metrics.Errors)
		fmt.Println("Status Success", metrics.StatusCodes)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
