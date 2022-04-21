package helpers

import (
	"log"
	"os/exec"
)

const (
	runStrykerCommand = "dotnet-stryker"
	strykerFileFlag   = "-f"
)

func RunStrykerMutator(configFile string) error {
	command := exec.Command(runStrykerCommand, strykerFileFlag, configFile)
	_, err := command.Output()
	command.Wait()
	if err != nil {
		log.Fatalf("Got error %v", err)
		return err
	}
	return nil
}
