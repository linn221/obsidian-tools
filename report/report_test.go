package report

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestGenerateUsesSectionNotename(t *testing.T) {
	report := &Report{
		Hot: []Section{{
			Date:     time.Date(2026, 11, 17, 0, 0, 0, 0, time.UTC),
			Notename: "2026-11-17",
			Headings: []string{"Hot"},
		}},
		Remind: []Section{{
			Date:     time.Date(2026, 11, 16, 0, 0, 0, 0, time.UTC),
			Notename: "2026-11-16",
			Headings: []string{"Remind"},
		}},
		MemorySheets: []Section{{
			Date:     time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC),
			Notename: "2026-11-15",
			Headings: []string{"Memory"},
		}},
	}

	var buf bytes.Buffer
	if err := report.Generate(&buf); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	output := buf.String()
	for _, want := range []string{
		"![[2026-11-17#Hot]]",
		"![[2026-11-16#Remind]]",
		"![[2026-11-15#Memory]]",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q\n%s", want, output)
		}
	}
}
