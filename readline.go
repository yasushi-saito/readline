// Package readline provides an easier-to-use wrapper around creadline.
//
// Example:
//   readline.Init(readline.Opts{Name: "myapp"})
//   for {
//     line = readline.Readline("> ")
//     processInput(line)
//     readline.AddHistory(line)
//   }
package readline

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/grailbio/base/errors"
	"github.com/yasushi-saito/readline/creadline"
)

// Opts define options for the readline package.
type Opts struct {
	// Name is the name of the application. It is used to generate the default
	// history pathname (~/.NAME_history).
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
	// last 10000 entries are retained.
	MaxHistoryLen int

	// ExpandHistory enables history expansion, such as "!tok".
	ExpandHistory bool

	// Completer is invoked to complete a line. It may be nil.
	Completer func(line string, start, end int) []string
}

var (
	opts          Opts
	curHistoryLen int
)

// Init sets the configuration parameter. It must be called at least once before
// calling any other function.
//
// Thread-hostile.
func Init(o Opts) error {
	opts = o
	err := errors.Once{}
	if opts.InitPath != "" {
		err.Set(creadline.ReadInitFile(opts.InitPath))
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
	err.Set(creadline.ReadHistory(opts.HistoryPath))
	if opts.Completer != nil {
		creadline.SetAttemptedCompletionFunction(opts.Completer)
	}
	if opts.MaxHistoryLen <= 0 {
		opts.MaxHistoryLen = 10000
	}
	creadline.StifleHistory(opts.MaxHistoryLen)
	curHistoryLen = creadline.HistoryLength()
	return err.Err()
}

// Readline reads one line. Thread-hostile.
func Readline(prompt string) (string, error) {
	for {
		line := creadline.Readline(prompt)
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
			return "", errors.E(line2)
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
