package main

import (
	"fmt"
	"github.com/xyproto/cdrom"
	"os"
)

func main() {
	cd, err := cdrom.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Ejecting")

	if err := cd.Eject(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := cd.Done(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
