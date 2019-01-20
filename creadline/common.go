package creadline

import (
	"errors"
)

// HistEntry is a Go translation of HIST_ENTRY
type HistEntry struct{ Line string }

// HistoryState is a Go translation of HISTORY_STATE
type HistoryState struct {
	Entries []HistEntry
	// Offset is the location pointer within Entries
	Offset int
	// Flags is a copy of HISTORY_STATE.flags.
	Flags int
}

// Interrupt is returned by Readline on SIGINT (i.e., Control-C keypress).
var Interrupt = errors.New("Interrupt")
