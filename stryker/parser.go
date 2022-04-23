package stryker

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"
)

type MutationReport struct {
	SchemaVersion string `json:"schemaVersion"`
	Thresholds    `json:"thresholds"`
	ProjectRoot   string                        `json:"projectRoot"`
	Files         map[string]MutationReportItem `json:"files"`
}

type Thresholds struct {
	High int `json:"high"`
	Low  int `json:"low"`
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

func ParseMutationReport(filePath string) MutationReport {
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
	// TODO: fix this
	for err == nil {
		if match := regex.MatchString(line); match {
			reportData := strings.Split(line, reportDataPrefix)[1]
			reportData = strings.TrimSuffix(reportData, reportDataSuffix)
			json.Unmarshal([]byte(reportData), &mutationReport)
			files := make(map[string]MutationReportItem)
			for key, value := range mutationReport.Files {
				testStrings := strings.Split(mutationReport.ProjectRoot, ".")
				finalString := testStrings[len(testStrings)-1]
				newKey := finalString + "\\" + key
				files[newKey] = value
			}
			mutationReport.Files = files
		}
		line, err = readLine(reader)
	}
	return mutationReport
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
