package stryker

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	reportDataPrefix      = "app.report = "
	reportDataSuffix      = ";"
	reportDataPrefixRegex = "^\\s*" + reportDataPrefix
)

func parseMutationReport(filePath string) MutationReport {
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
			mutationReport.Files = parseMutationFiles(mutationReport, line)
		}
		line, err = readLine(reader)
	}
	return mutationReport
}

func parseMutationFiles(mutationReport MutationReport, line string) map[string]MutationReportItem {
	files := make(map[string]MutationReportItem)
	for key, value := range mutationReport.Files {
		projectRootStrings := strings.Split(mutationReport.ProjectRoot, ".")
		lastProjectRootString := projectRootStrings[len(projectRootStrings)-1]
		newKey := lastProjectRootString + "\\" + key
		files[newKey] = value
	}
	return files
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
