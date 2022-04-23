package stryker

import (
	"bytes"
	"embed"
	"encoding/json"
	"io"
	"log"
	"text/template"
)

var (
	//go:embed "templates/*.gohtml"
	reportTemplates embed.FS
)

type ReportRenderer struct {
	templ *template.Template
}

func RenderHtmlReport(report MutationReport) string {
	reportRenderer, err := newReportRenderer()
	if err != nil {
		log.Fatalln(err)
	}

	buf := bytes.Buffer{}
	if err := reportRenderer.render(&buf, report); err != nil {
		log.Fatalln(err)
	}

	return buf.String()
}

func newReportRenderer() (*ReportRenderer, error) {
	templ, err := template.ParseFS(reportTemplates, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	return &ReportRenderer{templ: templ}, nil
}

type ReportString struct {
	Files string
}

func (r *ReportRenderer) render(w io.Writer, report MutationReport) error {
	jsonReportData, _ := json.Marshal(report)
	stringReportData := ReportString{string(jsonReportData)}
	if err := r.templ.Execute(w, stringReportData); err != nil {
		return err
	}

	return nil
}
