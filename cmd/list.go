/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.Default()

		accessToken, err := refreshToken()
		if err != nil {
			logger.Error("err refreshToken", "err", err)
			os.Exit(1)
		}

		op, err := openings(accessToken)
		if err != nil {
			logger.Error("err openings", "err", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(op, "", "  ")
		if err != nil {
			logger.Error("err json.MarshalIndent", "err", err)
			os.Exit(1)
		}

		fmt.Println(string(output))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
