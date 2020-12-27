package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/xyproto/cdrom"
	"github.com/xyproto/textoutput"
)

const versionString = "cupholder 1.1.0"

func ejectDevice(o *textoutput.TextOutput, deviceFilename string) error {
	o.Printf("<darkgray>[<blue>cupholder<darkgray>]\t\t<darkgray>Ejecting <yellow>%s<darkgray>... <off>", deviceFilename)

	// Opening
	cd, err := cdrom.NewFile(deviceFilename)
	if err != nil {
		o.Print("<red>error:</red> ")
		return err
	}

	// Ejecting
	if err := cd.Eject(); err != nil {
		o.Print("<red>error:</red> ")
		_ = cd.Done() // Try closing, but only return the other error
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

func generateEjectionHandler(deviceFilenames []string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "<!doctype html><html><head><title>Cupholder</title><style>html { height: 100%; } body { height: 100%; background: rgb(95,159,221); background-image: linear-gradient(rgba(95,159,221,1) 0%, rgba(95,159,221,1) 6em, rgba(218,192,168,1) 6em, rgba(224,250,255,1) 100%); margin: 3em; overflow: hidden; } h1 { color: orange; text-shadow: 2px 2px 4px #000F; }</style><body>")
		fmt.Fprintf(w, "<h1>%s</h1><br>", versionString)
		for _, deviceFilename := range deviceFilenames {
			fmt.Fprintf(w, "<strong>[cupholder]</strong>&nbsp;&nbsp;<code>Ejecting %s...</code> ", deviceFilename)

			// Opening
			cd, err := cdrom.NewFile(deviceFilename)
			if err != nil {
				fmt.Fprintf(w, "<code>error: %v</code><br>", err)
			}

			// Ejecting
			if err := cd.Eject(); err != nil {
				fmt.Fprintf(w, "<code>error: %v</code><br>", err)
				_ = cd.Done() // Try closing, but only output the other error
				continue
			}

			// Closing
			if err := cd.Done(); err != nil {
				fmt.Fprintf(w, "<code>error: %v</code><br>", err)
				continue
			}

			// All done, for this device filename
			fmt.Fprintln(w, "<code>ok</code><br>")
		}
	}
}

func main() {
	o := textoutput.New()
	if appErr := (&cli.App{
		Name:  "cupholder",
		Usage: "eject the CD tray (or other trays, given a device file)",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "silent", Aliases: []string{"s"}},
			&cli.BoolFlag{Name: "version", Aliases: []string{"V"}},
			&cli.BoolFlag{Name: "server", Aliases: []string{"l", "listen", "web", "http"}},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				o.Println(versionString)
				os.Exit(0)
			}
			// Check if text output should be disabled
			if c.Bool("silent") {
				o.Disable()
			}
			deviceFilenames := []string{"/dev/cdrom"}
			// Check if any arguments are given
			if c.NArg() > 0 {
				deviceFilenames = c.Args().Slice()
			}
			if c.Bool("server") {
				// Set up a HTTP server that is ready to eject the CD ROM at requests on port 0x0CD0 (3280)
				http.HandleFunc("/", generateEjectionHandler(deviceFilenames))
				o.Println("<blue>Listening for ejection requests at <yellow>http://localhost:3280/<off>")
				return http.ListenAndServe(":3280", nil)
			}
			// Treat all arguments as device files that shall be ejected
			var err error
			for _, deviceFilename := range deviceFilenames {

				err = ejectDevice(o, deviceFilename)
				if err != nil {
					o.Printf("<darkred>%s</darkred>\n", err)
					continue
				}
			}
			return err
		},
	}).Run(os.Args); appErr != nil {
		o.ErrExit(appErr.Error())
	}
}
