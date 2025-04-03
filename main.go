package main

import (
	"fmt"
	"os"

	"github.com/sharik709/bootstraper/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
