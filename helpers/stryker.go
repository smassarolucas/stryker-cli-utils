package helpers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
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

func MergeStrykerReports(filePaths []string) string {

	// TODO: Make use of goroutines to parse files separatedly
	for _, filePath := range filePaths {
		fmt.Println(filePath)
	}
	return ""
}

type MutationReport struct {
	SchemaVersion string `json:"schemaVersion"`
	Thresholds    struct {
		High int `json:"high"`
		Low  int `json:"low"`
	} `json:"thresholds"`
	ProjectRoot string                        `json:"projectRoot"`
	Files       map[string]MutationReportItem `json:"files"`
}

type MutationReportItem struct {
	Language string `json:"language"`
	Source   string `json:"source"`
	Mutants  []struct {
		ID          string `json:"id"`
		MutatorName string `json:"mutatorName"`
		Replacement string `json:"replacement"`
		Location    struct {
			Start struct {
				Line   int `json:"line"`
				Column int `json:"column"`
			} `json:"start"`
			End struct {
				Line   int `json:"line"`
				Column int `json:"column"`
			} `json:"end"`
		} `json:"location"`
		Status string `json:"status"`
		Static bool   `json:"static"`
	} `json:"mutants"`
}

const (
	reportDataPrefix      = "app.report = "
	reportDataSuffix      = ";"
	reportDataPrefixRegex = "^\\s*" + reportDataPrefix
)

func ParseMutationReport(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	regex, err := regexp.Compile(reportDataPrefixRegex)
	if err != nil {
		log.Fatalf("Couldn't compile the regex %v because of error %v", reportDataPrefixRegex, err.Error())
	}

	reader := bufio.NewReader(file)
	line, err := readLine(reader)
	mutationReport := MutationReport{}
	for err == nil {
		if match := regex.MatchString(line); match {
			reportData := strings.Split(line, reportDataPrefix)[1]
			reportData = strings.TrimSuffix(reportData, reportDataSuffix)
			json.Unmarshal([]byte(reportData), &mutationReport)
			// TODO: Alter the keys in the Files map to properly show the different projects

			// TODO: Make this another method, running only once when the goroutines finish their processing
			jsonReportData, _ := json.Marshal(mutationReport)
			stringReportData := string(jsonReportData)
			fmt.Println(string(stringReportData))
		}
		line, err = readLine(reader)
	}
}

func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix              bool  = true
		err                   error = nil
		currentLine, fullLine []byte
	)
	for isPrefix && err == nil {
		currentLine, isPrefix, err = r.ReadLine()
		fullLine = append(fullLine, currentLine...)
	}
	return string(fullLine), err
}
