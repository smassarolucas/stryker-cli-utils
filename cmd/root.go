/*
Copyright © 2022 Lucas Smassaro at smassarolucas@gmail.com

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stryker-cli-utils",
	Short: "Utils for improving Stryker .NET experience.",
	Long: `This package aims to be an improvement to the mutation testing experience with Stryker .NET (util now). Features, at the moment, include:
	
- Running all the mutation tests for the solutions and condensing the report
- Passing the report location via flag`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
