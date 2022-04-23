package stryker

import (
	"embed"
	"encoding/json"
	"io"
	"text/template"
)

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
