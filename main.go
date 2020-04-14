package main

import (
	"github.com/urfave/cli/v2"
	"github.com/xyproto/ask"
	"github.com/xyproto/clip"
	"github.com/xyproto/textoutput"
	"os"
)

const (
	versionString = "pastefile 0.4.0"
)

// Write to a file, using the contents from the clipboard
func writeFromClipboard(filename string) (int, error) {
	// Read the clipboard
	contents, err := clip.ReadAllBytes()
	if err != nil {
		return 0, err
	}
	// Write to file
	f, err := os.Create(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(contents)
}

func main() {
	o := textoutput.New()
	if appErr := (&cli.App{
		Name:  "pastefile",
		Usage: "create a file that contains the contents from the clipboard",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "silent", Aliases: []string{"s"}},
			&cli.BoolFlag{Name: "version", Aliases: []string{"V"}},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("version") {
				o.Println(versionString)
				os.Exit(0)
			}
			var filename string
			// Check if a filename is given
			if c.NArg() > 0 {
				filename = c.Args().Slice()[0]
			} else {
				o.Print("<yellow>filename:</yellow> <cyan>")
				filename = ask.ReadLn()
				o.Print("<cyan>")
			}
			bytesWritten, err := writeFromClipboard(filename)
			if err != nil {
				o.ErrExit(err.Error())
			}
			if !c.Bool("silent") {
				plural := ""
				if bytesWritten != 1 {
					plural = "s"
				}
				o.Printf("<blue>%d byte%s</blue> <off>written to <cyan>%s</cyan>\n", bytesWritten, plural, filename)
			}
			return nil
		},
	}).Run(os.Args); appErr != nil {
		o.ErrExit(appErr.Error())
	}
}
