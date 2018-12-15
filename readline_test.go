package readline_test

import (
	"log"
	"testing"

	"github.com/grailbio/testutil/expect"
	"github.com/yasushi-saito/readline"
)

func TestReadline(t *testing.T) {
	readline.Init(readline.Opts{
		Name:          "goreadlinetest",
		ExpandHistory: true,
	})
	for {
		line, err := readline.Readline("aueo>")
		log.Printf("Got: '%s' %v", line, err)
		expect.NoError(t, readline.AddHistory(line))
	}
}
