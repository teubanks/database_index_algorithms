package main

import (
	"fmt"
	"os"

	"example.com/learn/btrees/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
