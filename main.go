package main

import (
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/xyproto/ask"
	"github.com/xyproto/clip"
	"github.com/xyproto/textoutput"
)

const (
	versionString = "pastefile 0.5.1"
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

// exists checks if the given path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
				for {
					o.Print("<yellow>filename:</yellow> <blue>")
					filename = ask.ReadLn()
					o.Print("</blue>")
					if len(strings.TrimSpace(filename)) == 0 {
						continue
					}
					if !exists(filename) {
						break
					}
					o.Println("<red>file already exists</red>")
				}
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
				o.Printf("<green>%d byte%s</green> <white>written to</white> <blue>%s</blue>\n", bytesWritten, plural, filename)
			}
			return nil
		},
	}).Run(os.Args); appErr != nil {
		o.ErrExit(appErr.Error())
	}
}
