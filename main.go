package main

import (
	"github.com/atotto/clipboard"
	"github.com/urfave/cli/v2"
	"github.com/xyproto/textoutput"
	"os"
)

const versionString = "pastefile 0.1.0"

func main() {
	o := textoutput.New()
	if appErr := (&cli.App{
		Name:  "pastefile",
		Usage: "create a file that contains the contents from the clipboard",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "version", Aliases: []string{"V"}},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				o.Println(versionString)
				os.Exit(0)
			}
			filename := "pastefile.txt"
			// Check if a filename is given
			if c.NArg() > 0 {
				filename = c.Args().Slice()[0]
			}
			// Read the clipboard
			contents, err := clipboard.ReadAll()
			if err != nil {
				o.ErrExit(err.Error())
			}
			// Write to file
			f, err := os.Create(filename)
			if err != nil {
				o.ErrExit(err.Error())
			}
			f.WriteString(contents)
			f.Close()
			return nil
		},
	}).Run(os.Args); appErr != nil {
		o.ErrExit(appErr.Error())
	}
}
