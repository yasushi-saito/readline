package creadline

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
