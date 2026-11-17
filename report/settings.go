package report

import (
	"regexp"
	"time"

	"github.com/linn221/obsidian-tools/note"
	"github.com/linn221/obsidian-tools/utils"
)

// Edit these values to configure Vaultmate.
// All paths are relative to VaultDir unless they are absolute paths.

const (
	VaultDir      = "/home/linn/vault"
	DailyNotesDir = VaultDir
	OutputDir     = VaultDir

	// The exact text of the level-three heading, without the ### prefix.
	Hot = `hot`

	// Daily notes must be named YYYY-MM-DD.md.
	DailyNoteDateFormat = "2006-01-02"
)

// ReviewIntervals are calendar-day offsets from the original daily note date.
var ReviewIntervals = map[int]struct{}{
	1:   {},
	2:   {},
	4:   {},
	7:   {},
	14:  {},
	30:  {},
	60:  {},
	100: {},
	200: {},
	300: {},
	365: {},
}

// remind in ? days
var RemindSheetRe1 = regexp.MustCompile(`remind in (\d+) day`)

// remind at ?
var RemindSheetRe2 = regexp.MustCompile(`remind at (.*)`)

var MemorySheetHeading = regexp.MustCompile("(?i)(longterm|long term|memory)")

func NoteToBeReviewed(n *note.Note, today time.Time) bool {
	diff := utils.MinusDay(today, n.Date)
	_, ok := ReviewIntervals[diff]
	return ok
}
