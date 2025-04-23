package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var appName string

	rootCmd := &cobra.Command{
		Use:   "secrets",
		Short: "Inject your environment secrets into the curernt shell",
		Run:   CreateCommand(&appName),
	}

	rootCmd.Flags().StringVar(&appName, "app", "", "Application name is required")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
