/*
Copyright © 2022 Lucas Smassaro at smassarolucas@gmail.com

*/
package cmd

import (
	"log"

	"github.com/smassarolucas/stryker-cli-utils/stryker"
	"github.com/spf13/cobra"
)

// TODO: Delete StrykerReport folder
var mutateCmd = &cobra.Command{
	Use:   "mutate",
	Short: "Runs .NET for all configs and condenses the reports in one",
	Long:  `Runs Stryker.NET for ALL projects ending with stryker-config.json, outputs the result as a single HTML file and cleans the other reports.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Execution started!")
		stryker.RunStrykerForAllConfigs()
		log.Println("Execution finished!")
		reportLocation, _ := cmd.Flags().GetString("reportLocation")
		reportPath := stryker.GenerateReport(reportLocation)
		log.Printf("Report generated at: %v\n", reportPath)
	},
}

func init() {
	rootCmd.AddCommand(mutateCmd)

	mutateCmd.Flags().String("reportLocation", "", "Location for the report to be generated")
}
