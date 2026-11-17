package report

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"text/template"
	"time"
)

type Report struct {
	MemorySheets []Section
	Hot          []Section
	Remind       []Section
	Tasks        []Task
}

type Section struct {
	Date     time.Time
	Notename string
	Headings []string
}

type Task struct {
	Date     time.Time
	Notename string
	Tasks    []string
}

const reportTemplate = `
# Hot
{{range .Hot}}
- [[{{.Notename}}]]
{{$section := .}}
{{range .Headings}}
![[{{$section.Notename}}#{{.}}]]
{{end}}
---
{{end}}

# Reminded
{{range .Remind}}
- [[{{.Notename}}]]
{{$section := .}}
{{range .Headings}}
![[{{$section.Notename}}#{{.}}]]
{{end}}
---
{{end}}

# MemorySheets
{{range .MemorySheets}}
- [[{{.Notename}}]]
{{$section := .}}
{{range .Headings}}
![[{{$section.Notename}}#{{.}}]]
{{end}}
---
{{end}}

# Unfinished Tasks
{{range .Tasks}}
- [[{{.Notename}}]]
{{range .Tasks}}
	- {{.}}
{{end}}
---
{{end}}
`

func (report *Report) Generate(w io.Writer) error {
	// Sort copies so the original report remains unchanged.
	sort.Slice(report.Hot, func(i, j int) bool {
		return report.Hot[i].Date.After(report.Hot[j].Date)
	})

	sort.Slice(report.Remind, func(i, j int) bool {
		return report.Remind[i].Date.After(report.Remind[j].Date)
	})

	sort.Slice(report.MemorySheets, func(i, j int) bool {
		return report.MemorySheets[i].Date.After(report.MemorySheets[j].Date)
	})

	sort.Slice(report.Tasks, func(i, j int) bool {
		return report.Tasks[i].Date.After(report.Tasks[j].Date)
	})

	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return fmt.Errorf("parse report template: %w", err)
	}

	var output bytes.Buffer

	if err := tmpl.Execute(&output, report); err != nil {
		return fmt.Errorf("render report template: %w", err)
	}

	if _, err := w.Write(output.Bytes()); err != nil {
		return fmt.Errorf("write report: %w", err)
	}

	return nil
}
