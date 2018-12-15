// +build cgo

// Package creadline provides low-level cgo interface to GNU readline.  Most
// functions invoke C counterparts directly.
//
// See also github.com/yasushi-saito/readline. It provides an easier to use
// interface on top of creadline.
package creadline

/*
#cgo darwin CFLAGS: -I/usr/local/opt/readline/include
#cgo darwin LDFLAGS: -L/usr/local/opt/readline/lib
#cgo CFLAGS: -Wall
#cgo LDFLAGS: -lreadline

#include <stdio.h>
#include <stdlib.h>
#include <readline/readline.h>
#include <readline/history.h>

extern char* _go_readline_complete();
extern int _go_history_len();

*/
import "C"
import (
	"os"
	"os/signal"
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	C.rl_catch_sigwinch = 0
	c := make(chan os.Signal, 1)
	signal.Notify(c, unix.SIGWINCH)
	go func() {
		for _ = range c {
			C.rl_resize_terminal()
		}
	}()
}

func errnoToError(err C.int) error {
	if err == 0 {
		return nil
	}
	return unix.Errno(err)
}

func toCStringOrNil(s string) *C.char {
	if s == "" {
		return nil
	}
	return C.CString(s)
}

func freeOrNil(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}
	C.free(ptr)
}

// Readline calls readline.
func Readline(prompt string) string {
	cprompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(cprompt))
	cline := C.readline(cprompt)
	defer C.free(unsafe.Pointer(cline))
	line := C.GoString(cline)
	return line
}

// AddHistory calls readline add_history
func AddHistory(str string) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	C.add_history(cstr)
}

// WriteHistory calls readline read_history
func ReadHistory(path string) error {
	cpath := toCStringOrNil(path)
	defer freeOrNil(unsafe.Pointer(cpath))
	return errnoToError(C.read_history(cpath))
}

// WriteHistory calls readline write_history
func WriteHistory(path string) error {
	cpath := toCStringOrNil(path)
	defer freeOrNil(unsafe.Pointer(cpath))
	return errnoToError(C.write_history(cpath))
}

// AppendHistory calls readline append_history.
func AppendHistory(n int, path string) error {
	cpath := toCStringOrNil(path)
	defer freeOrNil(unsafe.Pointer(cpath))
	return errnoToError(C.append_history(C.int(n), cpath))
}

// HistoryTruncateFile calls readline history_truncate_file.
func HistoryTruncateFile(path string, n int) error {
	cpath := toCStringOrNil(path)
	defer freeOrNil(unsafe.Pointer(cpath))
	return errnoToError(C.history_truncate_file(cpath, C.int(n)))
}

// ClearHistory calls clear_history.
func ClearHistory() {
	C.clear_history()
}

// StifleHistory calls stifle_history
func StifleHistory(n int) {
	C.stifle_history(C.int(n))
}

// UnstifleHistory calls unstifle_history
func UnstifleHistory() int {
	return int(C.unstifle_history())
}

// HistoryLen reports the number of entries in the in-memory history list.
func HistoryLen() int {
	return int(C._go_history_len())
}

// ReadInitFile rl_calls read_init_file.
func ReadInitFile(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return errnoToError(C.rl_read_init_file(cpath))
}

// SetAttemptedCompletionFunction sets the completer function.
func SetAttemptedCompletionFunction(fn func(line string, start, end int) []string) {
	completionFunction = fn
	if fn == nil {
		C.rl_attempted_completion_function = nil
		return
	}
	C.rl_attempted_completion_function = (*C.rl_completion_func_t)(C._go_readline_complete)
}

var completionFunction func(line string, start, end int) []string

//export _goReadlineComplete
func _goReadlineComplete(_ *C.char, start, end C.int) **C.char {
	line := C.GoString(C.rl_line_buffer)
	completions := completionFunction(line, int(start), int(end))
	const ptrSize = unsafe.Sizeof((*C.char)(nil))
	array := C.malloc(C.size_t(len(completions)+1) * C.size_t(ptrSize))
	slot := func(i int) **C.char {
		return (**C.char)(unsafe.Pointer(uintptr(array) + uintptr(i)*ptrSize))
	}
	for i, completion := range completions {
		*slot(i) = (*C.char)(C.CString(completion))
	}
	*slot(len(completions)) = nil
	return (**C.char)(array)
}
