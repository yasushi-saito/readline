// +build cgo

package readline

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

func Readline(prompt string) string {
	cprompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(cprompt))
	cline := C.readline(cprompt)
	defer C.free(unsafe.Pointer(cline))
	line := C.GoString(cline)
	return line
}

func AddHistory(str string) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	C.add_history(cstr)
}

func ReadHistory(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return errnoToError(C.read_history(cpath))
}

func WriteHistory(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return errnoToError(C.write_history(cpath))
}

func AppendHistory(n int, path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return errnoToError(C.append_history(C.int(n), cpath))
}

func HistoryTruncateFile(path string, n int) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return errnoToError(C.history_truncate_file(cpath, C.int(n)))
}

func ReadInitFile(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return errnoToError(C.rl_read_init_file(cpath))
}

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
