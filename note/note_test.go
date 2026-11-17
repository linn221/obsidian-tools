package note

import (
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseNote(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		markdown      string
		wantDate      time.Time
		wantSections  []Header
		wantTasks     []string
		wantCompleted []string
	}{
		{
			name:     "basic markdown",
			filename: "2026-11-17.md",
			markdown: `# Daily Note

## Work

### Longterm

- [ ] Learn Go concurrency
- [x] Understand worker pools
- [ ] Read about database indexing

### Remind Nov 18

- Buy groceries

## Personal

- [ ] Read a book
`,
			wantDate: time.Date(2026, 11, 17, 0, 0, 0, 0, time.UTC),
			wantSections: []Header{
				{Name: "Daily Note", Type: H1},
				{Name: "Work", Type: H2},
				{Name: "Longterm", Type: H3},
				{Name: "Remind Nov 18", Type: H3},
				{Name: "Personal", Type: H2},
			},
			wantTasks: []string{
				"Learn Go concurrency",
				"Read about database indexing",
				"Read a book",
			},
			wantCompleted: []string{
				"Understand worker pools",
			},
		},
		{
			name:     "formatted task descriptions",
			filename: "2026-11-18.md",
			markdown: `# Test

- [ ] Learn **Go concurrency**
- [x] Read [Goldmark documentation](https://github.com/yuin/goldmark)
- [ ] Use ` + "`code`" + ` in a task
`,
			wantDate: time.Date(2026, 11, 18, 0, 0, 0, 0, time.UTC),
			wantSections: []Header{
				{Name: "Test", Type: H1},
			},
			wantTasks: []string{
				"Learn Go concurrency",
				"Use code in a task",
			},
			wantCompleted: []string{
				"Read Goldmark documentation",
			},
		},
		{
			name:     "no tasks",
			filename: "2026-11-19.md",
			markdown: `# Notes

## Work

Some ordinary text.
`,
			wantDate: time.Date(2026, 11, 19, 0, 0, 0, 0, time.UTC),
			wantSections: []Header{
				{Name: "Notes", Type: H1},
				{Name: "Work", Type: H2},
			},
		},
		{
			name:     "empty markdown",
			filename: "2026-11-20.md",
			markdown: "",
			wantDate: time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "invalid date filename",
			filename: "not-a-date.md",
			markdown: "# Notes\n",
			wantSections: []Header{
				{Name: "Notes", Type: H1},
			},
		},
		{
			name:     "nested filepath",
			filename: "Daily Notes/2026-11-21.md",
			markdown: "# Nested\n",
			wantDate: time.Date(2026, 11, 21, 0, 0, 0, 0, time.UTC),
			wantSections: []Header{
				{Name: "Nested", Type: H1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNote(strings.NewReader(tt.markdown), tt.filename)
			if err != nil {
				t.Fatalf("ParseNote() error = %v", err)
			}

			if got.Filepath != tt.filename {
				t.Errorf("Filepath = %q, want %q", got.Filepath, tt.filename)
			}

			wantName := strings.TrimSuffix(filepath.Base(tt.filename), ".md")
			if got.Notename != wantName {
				t.Errorf("Notename = %q, want %q", got.Notename, wantName)
			}

			if got.Raw != tt.markdown {
				t.Errorf("Raw = %q, want %q", got.Raw, tt.markdown)
			}

			if !got.Date.Equal(tt.wantDate) {
				t.Errorf("Date = %v, want %v", got.Date, tt.wantDate)
			}

			if !equalHeaders(got.Sections, tt.wantSections) {
				t.Errorf("Sections = %#v, want %#v", got.Sections, tt.wantSections)
			}

			if !equalStrings(got.Tasks, tt.wantTasks) {
				t.Errorf("Tasks = %#v, want %#v", got.Tasks, tt.wantTasks)
			}

			if !equalStrings(got.CompletedTasks, tt.wantCompleted) {
				t.Errorf("CompletedTasks = %#v, want %#v",
					got.CompletedTasks, tt.wantCompleted)
			}
		})
	}
}

func TestParseNoteReadError(t *testing.T) {
	reader := &errorReader{}

	_, err := ParseNote(reader, "2026-11-17.md")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

type errorReader struct{}

func (*errorReader) Read([]byte) (int, error) {
	return 0, errors.New("simulated read failure")
}

func equalHeaders(a, b []Header) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

var _ io.Reader = (*errorReader)(nil)
