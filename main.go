package main

import (
	"github.com/urfave/cli/v2"
	"github.com/xyproto/cdrom"
	"github.com/xyproto/textoutput"
	"os"
)

func ejectDevice(o *textoutput.TextOutput, deviceFilename string) error {
	o.Printf("<darkgray>[<blue>eject<darkgray>]\t\t<darkgray>Ejecting <yellow>%s<darkgray>... <off>", deviceFilename)

	// Opening
	cd, err := cdrom.NewFile(deviceFilename)
	if err != nil {
		o.Print("<red>error:</red> ")
		return err
	}

	// Ejecting
	if err := cd.Eject(); err != nil {
		o.Print("<red>error:</red> ")
		return err
	}

	// Closing
	if err := cd.Done(); err != nil {
		o.Print("<red>error:</red> ")
		return err
	}

	// All done
	o.Println("<green>ok</green>")
	return nil
}

func main() {
	o := textoutput.New()
	if appErr := (&cli.App{
		Name:  "eject",
		Usage: "eject the CD tray (or other trays, given a device file)",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "silent", Aliases: []string{"s"}},
		},
		Action: func(c *cli.Context) error {
			// Check if text output should be disabled
			if c.Bool("silent") {
				o.Disable()
			}
			// Check if any arguments are given
			if c.NArg() == 0 {
				return ejectDevice(o, "/dev/cdrom")
			}
			// Treat all arguments as device files that shall be ejected
			var err error
			for _, deviceFilename := range c.Args().Slice() {
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
