package readline_test

import (
	"log"
	"testing"

	"github.com/grailbio/testutil/expect"
	"github.com/yasushi-saito/readline"
)

func TestReadline(t *testing.T) {
	readline.Init(readline.Opts{
		Name: "goreadlinetest",
	})
	for {
		line := readline.Readline("aueo>")
		log.Printf("Got: '%s'", line)
		expect.NoError(t, readline.AddHistory(line))
	}
}
