/*
Copyright © 2022 Lucas Smassaro at smassarolucas@gmail.com

*/
package cmd

import (
	"log"

	"github.com/smassarolucas/stryker-cli-utils/stryker"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs Stryker.NET for multiple projects and merges the reports into one",
	Long:  `Runs Stryker.NET for multiple projects, outputs the result as a single HTML file and cleans the other reports.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Execution started!")

		individualConfigsFlag, _ := cmd.Flags().GetBool("i")
		if individualConfigsFlag {
			log.Println("Running Stryker for each individual config file.")
			stryker.RunStrykerForAllConfigs()
		}

		baseConfigFlag, _ := cmd.Flags().GetBool("a")
		if baseConfigFlag {
			log.Println("Running Stryker for all projects referenced by test project.")
			stryker.RunStrykerForAllProjects()
		}

		log.Println("Execution finished!")
		reportLocation, _ := cmd.Flags().GetString("reportLocation")
		reportPath := stryker.GenerateReport(reportLocation)
		log.Printf("Report generated at: %v\n", reportPath)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().Bool("i", false, "This flag informs the runner that the execution will happen based on individual config files, suffixed by \"-stryker-config.json\"")
	runCmd.Flags().Bool("a", false, "This flag informs the runner that the execution will happen based on a base config file, getting the references from the project automatically")
	runCmd.Flags().String("reportLocation", "", "Location for the report to be generated. IMPORTANT: should be used inside double quotes like \"C:/dev/report.html\"")
}
