package stryker

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func RunStrykerForAllConfigs() string {
	fileNames := GetStrykerConfigFileNames()
	for _, fileName := range fileNames {
		runStrykerMutator(fileName)
	}
	filePaths := GetMutationReportsFilePaths()
	return MergeStrykerReports(filePaths)
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

func MergeStrykerReports(filePaths []string) string {
	strykerReports := make([]MutationReport, len(filePaths))
	// TODO: Make use of goroutines to parse files separatedly
	// TODO: Fix the threshold thing
	thresholds := Thresholds{High: 0, Low: 0}
	for _, filePath := range filePaths {
		strykerReport := ParseMutationReport(filePath)
		thresholds.High += strykerReport.Thresholds.High
		thresholds.Low += strykerReport.Thresholds.Low
		strykerReports = append(strykerReports, strykerReport)
	}
	thresholds.High = thresholds.High / len(strykerReports)
	thresholds.Low = thresholds.Low / len(strykerReports)

	reportFiles := make(map[string]MutationReportItem)
	for _, report := range strykerReports {
		for key, value := range report.Files {
			reportFiles[key] = value
		}
	}

	finalReport := MutationReport{
		SchemaVersion: strykerReports[0].SchemaVersion,
		Thresholds:    thresholds,
		ProjectRoot:   strykerReports[0].ProjectRoot,
		Files:         reportFiles,
	}

	// TODO: Segregate this concern
	reportRenderer, err := NewReportRenderer()
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.Buffer{}
	if err := reportRenderer.Render(&buf, finalReport); err != nil {
		log.Fatalln(err)
	}

	// TODO: Segregate this concern
	f, err := os.Create("./consolidated-report.html")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	f.WriteString(buf.String())

	filePath, _ := filepath.Abs(f.Name())

	return filePath
}
