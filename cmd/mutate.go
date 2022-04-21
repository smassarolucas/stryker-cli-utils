/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/smassarolucas/stryker-cli-utils/helpers"
	"github.com/spf13/cobra"
)

// mutateCmd represents the mutate command
var mutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Runs .NET for all configs and condenses the reports in one",
	Long:  `Runs Stryker.NET for ALL projects ending with stryker-config.json, outputs the result as a single HTML file and cleans the other reports.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mutate called")
		fileNames := helpers.GetStrykerConfigFileNames()
		// TODO: Run Stryker.NET for all configs
		for _, fileName := range fileNames {
			helpers.RunStrykerMutator(fileName)
		}

		// TODO: Merge the reports and output the location
	},
}

func init() {
	rootCmd.AddCommand(mutateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mutateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mutateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
