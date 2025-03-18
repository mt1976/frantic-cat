package reporthandler

import (
	"context"
	"strings"
	"time"

	"github.com/mt1976/frantic-cat/app/dao/report"
	"github.com/mt1976/frantic-core/idHelpers"
)

type Report struct {
	content    report.Report_Store
	reportType ReportType
}

var TYPE_Markdown = ReportType{Title: "Markdown"}
var TYPE_Default = ReportType{Title: "Default"}

type ReportType struct {
	Title string
}

func NewReport(name string, reportType ReportType) (Report, error) {
	t := time.Now()
	r := Report{}
	r.reportType = reportType
	r.content.Title = name
	r.content.Generated = t
	r.content.Host = "localhost"
	r.content.HostIP = ""
	r.content.Raw = name + "-" + t.Format("20060102150405") + "-" + r.content.Host
	r.content.Key = idHelpers.Encode(r.content.Raw)
	return r, nil
}

func (r *Report) AddRow(text string) {
	index := len(r.content.Content) + 1
	r.content.Content = append(r.content.Content, report.Row{Index: index, Text: text})
}

func (r *Report) Break() {
	r.AddRow(" ")
}

func (r *Report) HR() {
	r.AddRow("-------------------------------------------------")
}

func (r *Report) H1(name string) {
	l := len(name)
	rule := strings.Repeat("-", l)
	r.Break()
	r.AddRow(rule)
	r.AddRow(name)
	r.AddRow(rule)
}

func (r *Report) Spool() error {
	err := r.content.Create(context.TODO(), "")
	if err != nil {
		return err
	}
	return nil
}
