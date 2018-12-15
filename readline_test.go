package readline_test

import (
	"testing"
	"log"
	"github.com/yasushi-saito/readline"
)

func TestReadline(t*testing.T) {
	line := readline.Readline("aueo>")
	log.Printf("Got: '%s'", line)
}
