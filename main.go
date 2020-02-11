package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xyproto/cdrom"
	"github.com/xyproto/textoutput"
	"os"
)

func ejectDevice(o *textoutput.TextOutput, deviceFilename string) error {
	cd, err := cdrom.NewFile(deviceFilename)
	if err != nil {
		return err
	}
	fmt.Print(o.LightTags(fmt.Sprintf(
		"<blue>Ejecting</blue> <lightyellow>%s</lightyellow><blue>...</blue>",
		deviceFilename)))
	if err := cd.Eject(); err != nil {
		o.OutputTags("<red>failed</red>")
		return err
	}
	if err := cd.Done(); err != nil {
		o.OutputTags("<red>failed</red>")
		return err
	}
	o.OutputTags("<green>ok</green>")
	return nil
}

func main() {
	o := textoutput.New()
	// quiet output: o.Disable()
	if appErr := (&cli.App{
		Name:  "eject",
		Usage: "eject the CD tray (or other trays, given a device file)",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return ejectDevice(o, "/dev/cdrom")
			}
			var err error
			for _, deviceFilename := range c.Args().Slice() {
				fmt.Println(deviceFilename)
				err = ejectDevice(o, deviceFilename)
				if err != nil {
					break
				}
			}
			return err
		},
	}).Run(os.Args); appErr != nil {
		o.ErrExit(appErr.Error())
	}
}
