package main

import (
	"os"

	"github.com/shota-imoto/tohyo"
)

func main() {
	err := tohyo.InteractiveVote(os.Stdin)
	if err != nil {
		panic(err)
	}
}
