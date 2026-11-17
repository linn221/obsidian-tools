package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/linn221/obsidian-tools/note"
	"github.com/linn221/obsidian-tools/report"
	"github.com/linn221/obsidian-tools/utils"
)

func main() {

	dir := report.VaultDir
	notes := make([]note.Note, 0)

	outputFilepath := filepath.Join(report.OutputDir, "output.md")
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		filename := filepath.Base(path)
		ext := filepath.Ext(filename)
		if info.Mode().IsRegular() && ext == ".md" && path != outputFilepath {
			func() {
				file, err := os.Open(path)
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				n, err := note.ParseNote(file, path)
				if err != nil {
					log.Fatal(err)
				}
				notes = append(notes, n)
			}()
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	var rep report.Report
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	for _, note := range notes {
		hotHeadings := make([]string, 0)
		remindHeadings := make([]string, 0)
		memoryHeadings := make([]string, 0)
		toBeReviewed := report.NoteToBeReviewed(&note, today)
		var tasks []string

		for _, h := range note.Sections {
			// memory sheet
			if toBeReviewed {
				if utils.IsRegexMatch(report.MemorySheetHeading, h.Name) {
					memoryHeadings = append(memoryHeadings, h.Name)
				}
			}
			// hot
			if strEqual(h.Name, report.Hot) {
				hotHeadings = append(hotHeadings, strings.TrimSpace(h.Name))
			}

			// remind
			utils.MatchRegex(report.RemindSheetRe1, h.Name, func(match string) {
				d, _ := strconv.Atoi(match)
				if d > 0 {
					remindDate := note.Date.AddDate(0, 0, d)
					if utils.DateEqual(remindDate, today) {
						remindHeadings = append(remindHeadings, h.Name)
					}
				}
			})
			utils.MatchRegex(report.RemindSheetRe2, h.Name, func(match string) {
				remindDate := utils.ParseMyDate(match)
				// could not parse the date, could be an issue, including it
				if remindDate.IsZero() {
					remindHeadings = append(remindHeadings, h.Name)
				}

				if utils.DateEqual(remindDate, today) {
					remindHeadings = append(remindHeadings, h.Name)
				}
			})
		}
		tasks = note.Tasks
		if len(memoryHeadings) > 0 {
			rep.MemorySheets = append(rep.MemorySheets, report.Section{
				Date:     note.Date,
				Notename: note.Notename,
				Headings: memoryHeadings,
			})
		}
		if len(hotHeadings) > 0 {
			rep.Hot = append(rep.Hot, report.Section{
				Date:     note.Date,
				Notename: note.Notename,
				Headings: hotHeadings,
			})
		}
		if len(remindHeadings) > 0 {
			rep.Remind = append(rep.Remind, report.Section{
				Date:     note.Date,
				Notename: note.Notename,
				Headings: remindHeadings,
			})
		}
		if len(tasks) > 0 {
			rep.Tasks = append(rep.Tasks, report.Task{
				Date:     note.Date,
				Notename: note.Notename,
				Tasks:    tasks,
			})
		}
	}

	file, err := os.OpenFile(outputFilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Println(today)
	err = rep.Generate(file)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func strEqual(a, b string) bool {
	return strings.EqualFold(a, b)
}
