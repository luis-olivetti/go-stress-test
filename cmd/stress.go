package cmd

import (
	"errors"

	"github.com/luis-olivetti/go-stresstest/internal"
	"github.com/spf13/cobra"
)

var stressCmd = &cobra.Command{
	Use:   "stress",
	Short: "Execute a stress test on the server.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	stressCmd.Flags().StringP("url", "u", "", "The URL to stress test.")
	stressCmd.Flags().IntP("concurrency", "c", 1, "The number of concurrent requests to make.")
	stressCmd.Flags().IntP("requests", "r", 1, "The number of requests to make.")
}
