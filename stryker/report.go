package stryker

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
	defaultReportName = "./consolidated-report.html"
)

func mergeStrykerReports(filePaths []string, reportLocation string) string {
	reports, thresholds := joinReports(filePaths)

	reportFiles := make(map[string]MutationReportItem)
	for _, report := range reports {
		for key, value := range report.Files {
			reportFiles[key] = value
		}
	}

	mergedReport := MutationReport{
		SchemaVersion: reports[0].SchemaVersion,
		Thresholds: Thresholds{
			High: thresholds.High / len(reports),
			Low:  thresholds.Low / len(reports),
		},
		ProjectRoot: reports[0].ProjectRoot,
		Files:       reportFiles,
	}

	report := renderHtmlReport(mergedReport)

	if reportLocation == "" {
		reportLocation = defaultReportName
	}
	filePath := writeToFile(report, reportLocation)

	return filePath
}

func joinReports(filePaths []string) ([]MutationReport, Thresholds) {
	strykerReports := make([]MutationReport, 0)
	thresholds := Thresholds{High: 0, Low: 0}
	for _, filePath := range filePaths {
		strykerReport := parseMutationReport(filePath)
		thresholds.High += strykerReport.Thresholds.High
		thresholds.Low += strykerReport.Thresholds.Low
		strykerReports = append(strykerReports, strykerReport)
	}
	return strykerReports, thresholds
}
