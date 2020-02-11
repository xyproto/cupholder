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
	cd.Eject()
	if err := cd.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("DONE")
}
