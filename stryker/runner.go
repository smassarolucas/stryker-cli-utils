package stryker

import (
	"log"
	"os/exec"
)

const (
	runStrykerCommand = "dotnet-stryker"
	strykerFileFlag   = "-f"
)

func RunStrykerForAllConfigs() {
	fileNames, err := getStrykerConfigFileNames()

	if err != nil {
		log.Fatalf("Couldn't get config files because of error: %v", err.Error())
	}
	for _, fileName := range fileNames {
		runStrykerMutatorForConfig(fileName)
	}
}

func RunStrykerForAllProjects() {
	mutableProjects, testProject := getProjectsToMutate()
	for _, mutableProject := range mutableProjects {
		runStrykerForProject(mutableProject, testProject)
	}
}

func runStrykerMutatorForConfig(configFile string) error {
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

const (
	strykerProjectFlag     = "-p"
	strykerTestProjectFlag = "-tp"
)

func runStrykerForProject(mutableProject, testProject string) error {
	log.Printf("Running Stryker for project %v", mutableProject)
	command := exec.Command(runStrykerCommand, strykerProjectFlag, mutableProject, strykerTestProjectFlag, testProject)
	_, err := command.Output()
	command.Wait()
	if err != nil {
		log.Fatalf("Got error %v", err)
		return err
	}
	return nil
}

func GenerateReport(reportLocation string) string {
	filePaths := getMutationReportsFilePaths()
	reportPath := mergeStrykerReports(filePaths, reportLocation)
	deleteStrykerOutputFolder()
	return reportPath
}
