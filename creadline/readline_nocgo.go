// +build !cgo

package creadline

import (
	"bufio"
	"os"
	"strings"
)

var stdinReader *bufio.Reader

func Readline(prompt string) (string, error) {
	if stdinReader == nil {
		stdinReader = bufio.NewReader(os.Stdin)
	}
	os.Stdout.Write([]byte(prompt))
	line, err := stdinReader.ReadString('\n')
	return strings.TrimSpace(line), err
}

func AddHistory(str string)                                                        {}
func ReadHistory(path string) error                                                { return nil }
func WriteHistory(path string) error                                               { return nil }
func AppendHistory(n int, path string) error                                       { return nil }
func HistoryTruncateFile(path string, n int) error                                 { return nil }
func ClearHistory()                                                                {}
func StifleHistory(n int)                                                          {}
func UnstifleHistory() int                                                         { return 0 }
func HistoryLength() int                                                           { return 0 }
func HistoryGetHistoryState() (h HistoryState)                                     { return }
func HistoryExpand(val string) (out string, ret int)                               { return val, 0 }
func ReadInitFile(path string) error                                               { return nil }
func SetAttemptedCompletionFunction(fn func(line string, start, end int) []string) { return }

func GetScreenSize() (int, int) { return 0, 0 }

func Init() {}
