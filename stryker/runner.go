package stryker

import (
	"log"
	"os/exec"
)

func RunStrykerForAllConfigs() {
	fileNames := GetStrykerConfigFileNames()
	for _, fileName := range fileNames {
		runStrykerMutator(fileName)
	}
}

func GenerateReport(reportLocation string) string {
	filePaths := GetMutationReportsFilePaths()
	reportPath := MergeStrykerReports(filePaths, reportLocation)
	return reportPath
}

const (
	runStrykerCommand = "dotnet-stryker"
	strykerFileFlag   = "-f"
)

func runStrykerMutator(configFile string) error {
	log.Printf("Running Stryker for config file %v", configFile)
	command := exec.Command(runStrykerCommand, strykerFileFlag, configFile)
	_, err := command.Output()
	command.Wait()
	if err != nil {
		log.Fatalf("Got error %v", err)
		return err
	}
	return nil
}
