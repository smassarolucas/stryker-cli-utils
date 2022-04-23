/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/smassarolucas/stryker-cli-utils/stryker"
	"github.com/spf13/cobra"
)

// mutateCmd represents the mutate command
// TODO: See what can be parallel
// TODO: Create tests
// TODO: Separate files and concerns
// TODO: Make the report location passable as a flag
// TODO: Delete StrykerReport folder
// TODO: Make sure there are no fixed folder things
// TODO: Remove all jank
// TODO: Log everything
var mutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Runs .NET for all configs and condenses the reports in one",
	Long:  `Runs Stryker.NET for ALL projects ending with stryker-config.json, outputs the result as a single HTML file and cleans the other reports.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Execution started!")
		reportPath := stryker.RunStrykerForAllConfigs()
		log.Println("Execution finished!")
		log.Printf("Report generated at: %v\n", reportPath)
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
