package main

import (
	"fmt"
	"os"

	"collector-go/cmd/root"
)

func main() {

	if err := root.GetRootCommand().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
