package readline_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/yasushi-saito/readline"
)

func TestReadline(t *testing.T) {
	readline.Init(readline.Opts{
		Name:          "goreadlinetest",
		ExpandHistory: true,
	})
	nRow, nCol := readline.GetScreenSize()
	log.Printf("Screen size: %d %d", nRow, nCol)
	n := 0
	for {
		line, err := readline.Readline(fmt.Sprintf("test%02d> ", n))
		n++
		fmt.Printf("Got: '%s' %v\n", line, err)
		if err := readline.AddHistory(line); err != nil {
			t.Fatal(err)
		}
	}
}
