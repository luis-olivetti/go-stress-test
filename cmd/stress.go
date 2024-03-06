package cmd

import (
	"errors"

	"github.com/luis-olivetti/go-stresstest/internal"
	"github.com/spf13/cobra"
)

var stressCmd = &cobra.Command{
	Use:   "stress",
	Short: "Execute a stress test on the server.",
	Long: `The stress command initiates a stress test on the specified server URL by simulating multiple concurrent requests.

This command sends a configurable number of requests to the specified URL concurrently, allowing you to stress-test the server's performance under load.

Example:
  stress -u http://example.com -c 10 -r 100

This example initiates a stress test on http://example.com with 10 concurrent requests and 100 total requests.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("url") {
			return errors.New("required flag 'url' not set")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		urlStr, _ := cmd.Flags().GetString("url")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		requests, _ := cmd.Flags().GetInt("requests")

		internal.NewStress(urlStr, concurrency, requests).Run()
	},
}

func init() {
	rootCmd.AddCommand(stressCmd)

	stressCmd.Flags().StringP("url", "u", "", "The URL to stress test. (required)")
	stressCmd.Flags().IntP("concurrency", "c", 1, "The number of concurrent requests to make.")
	stressCmd.Flags().IntP("requests", "r", 1, "The number of requests to make.")
}
