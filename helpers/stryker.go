package helpers

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
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

var (
	//go:embed "templates/*.gohtml"
	reportTemplates embed.FS
)

type ReportRenderer struct {
	templ *template.Template
}

func NewReportRenderer() (*ReportRenderer, error) {
	templ, err := template.ParseFS(reportTemplates, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	return &ReportRenderer{templ: templ}, nil
}

type ReportString struct {
	Files string
}

func (r *ReportRenderer) Render(w io.Writer, report MutationReport) error {
	jsonReportData, _ := json.Marshal(report)
	stringReportData := ReportString{string(jsonReportData)}
	if err := r.templ.Execute(w, stringReportData); err != nil {
		return err
	}

	return nil
}

func MergeStrykerReports(filePaths []string) string {
	strykerReports := make([]MutationReport, len(filePaths))

	// TODO: Make use of goroutines to parse files separatedly
	for _, filePath := range filePaths {
		strykerReport := ParseMutationReport(filePath)
		strykerReports = append(strykerReports, strykerReport)
	}

	reportFiles := make(map[string]MutationReportItem)
	for _, report := range strykerReports {
		for key, value := range report.Files {
			reportFiles[key] = value
		}
	}

	finalReport := MutationReport{
		SchemaVersion: strykerReports[0].SchemaVersion,
		Thresholds:    strykerReports[0].Thresholds,
		ProjectRoot:   strykerReports[0].ProjectRoot,
		Files:         reportFiles,
	}

	reportRenderer, err := NewReportRenderer()
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.Buffer{}
	if err := reportRenderer.Render(&buf, finalReport); err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create("./test.html")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	f.WriteString(buf.String())

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
