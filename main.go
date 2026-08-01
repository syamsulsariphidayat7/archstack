package main

import (
	"fmt"
	"os"

	"github.com/syamsulsariphidayat7/archstack/internal/cli"
)

func main() {
	if err := cli.Root(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
