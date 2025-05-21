/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ssdSSF/swing/pkg/model"
	"gopkg.in/yaml.v2"
)

var cmdSecretsPath string
var cmdSecrets model.Secrets

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmdSecretsPath = expandPath(cmdSecretsPath)

		_, err := os.Stat(cmdSecretsPath)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "File %s not found", cmdSecretsPath)
			os.Exit(1)
		}

		f, err := os.Open(cmdSecretsPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err os.Open(cmdSecretPath)")
			os.Exit(1)
		}

		err = yaml.NewDecoder(f).Decode(&cmdSecrets)
		if err != nil {
			fmt.Fprintf(os.Stderr, "err json.NewDecoder(f).Decode(&secrets)")
			os.Exit(1)
		}

		if cmdSecrets.GoogleToken == "" {
			fmt.Fprintf(os.Stderr, "err secrets.GoogleToken is empty")
			os.Exit(1)
		}

		if cmdSecrets.SlackToken == "" {
			fmt.Fprintf(os.Stderr, "err secrets.SlackToken is empty")
			os.Exit(1)
		}

		if cmdSecrets.SlackChannel == "" {
			fmt.Fprintf(os.Stderr, "err secrets.SlackChannel is empty")
			os.Exit(1)
		}

		// fastest polling, 10 seconds
		if cmdSecrets.Interval < 10 {
			cmdSecrets.Interval = 10
		}

		if cmdSecrets.SlackHeartbeatChannel == "" {
			fmt.Fprintf(os.Stderr, "err secrets.SlackHeartbeatChannel is empty")
			os.Exit(1)
		}

		for i, city := range cmdSecrets.CitiesToSkip {
			cmdSecrets.CitiesToSkip[i] = strings.ToLower(city)
		}
	},
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path // Return original path on error
		}
		return filepath.Join(homeDir, path[1:])
	}
	return path
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&cmdSecretsPath, "secret-path", "p", "~/.swing-secrets.yaml", "Path of the Swing secrets including Google token and Slack token")
}
