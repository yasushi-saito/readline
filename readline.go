// Package readline provides an easier-to-use wrapper around creadline.
//
// Example:
//   readline.Init(readline.Opts{Name: "myapp"})
//   for {
//     line = readline.Readline("> ")
//     processInput(line)
//     readline.AddHistory(line)
//   }
//
// Signal handling
//
// Readline intercepts SIGINT (Ctrl-C). On SIGINT, the ongoing Readline call
// aborts and returns an Interrupt error.
package readline

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/yasushi-saito/readline/creadline"
)

// Opts define options for the readline package.
type Opts struct {
	// Name is the name of the application. It is used to generate the default
	// history pathname (~/.NAME_history). It may be empty.
	Name string
	// InitPath is the name of the readline init file. It may be empty.
	//
	// https://tiswww.case.edu/php/chet/readline/readline.html#SEC9
	InitPath string
	// HistoryPath is the name of the readline history file. If "" and Name is
	// nonempty, then its value will be ~/.NAME_history. If both HistoryPath and
	// Name are empty, then its value will be ~/.history.
	HistoryPath string
	// MaxHistoryLen is the maximum number of history entries to retain. If <= 0,
	// the last 10000 entries will be retained.
	MaxHistoryLen int
	// ExpandHistory enables history expansion, such as "!tok".
	ExpandHistory bool

	// Completer is invoked to complete a line. It may be nil.  The arg line is
	// the current input line. Args start and end are the start and limit offset
	// of the word being completed, respectively.
	Completer func(line string, start, end int) []string
}

var (
	// Interrupt is returned by Readline on SIGINT (i.e., Control-C keypress).
	Interrupt = creadline.Interrupt

	opts          Opts
	curHistoryLen int
)

// Init sets the configuration parameter. It must be called at least once before
// calling any other function.
//
// Thread-hostile.
func Init(o Opts) error {
	creadline.Init()
	opts = o
	if opts.InitPath != "" {
		if err := creadline.ReadInitFile(opts.InitPath); err != nil {
			return err
		}
	}
	if opts.HistoryPath == "" {
		usr, err := user.Current()
		if err != nil {
			return err
		}
		if opts.Name != "" {
			opts.HistoryPath = filepath.Join(usr.HomeDir, "."+opts.Name+"_history")
		} else {
			opts.HistoryPath = filepath.Join(usr.HomeDir, ".history")
		}
	}
	if err := creadline.ReadHistory(opts.HistoryPath); err != nil {
		return err
	}
	if opts.Completer != nil {
		creadline.SetAttemptedCompletionFunction(opts.Completer)
	}
	if opts.MaxHistoryLen <= 0 {
		opts.MaxHistoryLen = 10000
	}
	creadline.StifleHistory(opts.MaxHistoryLen)
	curHistoryLen = creadline.HistoryLength()
	return nil
}

// Readline reads one line. On SIGINT, it returns an Inturrupt error.
// Thread-hostile.
//
// REQUIRES: Init has been called.
func Readline(prompt string) (string, error) {
	for {
		line, err := creadline.Readline(prompt)
		if err != nil {
			return line, err
		}
		if !opts.ExpandHistory {
			return line, nil
		}
		line2, ret := creadline.HistoryExpand(line)
		switch ret {
		case 0:
			return line, nil
		case 1:
			return line2, nil
		case -1:
			return "", fmt.Errorf("history: %s", line2)
		case 2:
			fmt.Fprintf(os.Stderr, "%s: %s\n", opts.Name, line2)
			continue
		default:
			panic(ret)
		}
	}
}

// AddHistory adds a history entry. It appends the entry both inmemory list and
// Opts.HistoryPath. It may truncate the history list if it grows larger than
// Opts.MaxHistoryLen. Thread-hostile.
//
// REQUIRES: Init has been called.
func AddHistory(str string) (err error) {
	if opts.HistoryPath == "" {
		panic("readline.Init not yet called")
	}
	creadline.AddHistory(str)
	if _, e := os.Stat(opts.HistoryPath); e != nil {
		err = creadline.WriteHistory(opts.HistoryPath)
	} else {
		err = creadline.AppendHistory(1, opts.HistoryPath)
	}
	curHistoryLen++
	if curHistoryLen >= 10000 && curHistoryLen >= opts.MaxHistoryLen*4 {
		creadline.HistoryTruncateFile(opts.HistoryPath, opts.MaxHistoryLen) // nolint: errcheck
		curHistoryLen = opts.MaxHistoryLen
	}
	return err
}

// GetScreenSize returns the current screen size, (#rows, #cols).  It returns
// nonpositive values when the stdout is not a terminal, or on any error.
func GetScreenSize() (int, int) {
	return creadline.GetScreenSize()
}
