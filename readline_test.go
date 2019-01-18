package readline_test

import (
	"log"
	"testing"

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
		if err := readline.AddHistory(line); err != nil {
			t.Fatal(err)
		}
	}
}
