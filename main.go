package main

import (
	"fmt"
	"os"
	"path"

	"github.com/retcon85/retcon-util-audio/cli"
	"github.com/retcon85/retcon-util-audio/model/psg"
)

// this will receive a value from the linker
var version string
var buildDate string

func main() {
	prog := path.Base(os.Args[0])
	app := cli.NewCli(prog)
	if version != "" {
		app.Version = func() string {
			return fmt.Sprintf("%s version %s\n\nbuild date: %s", prog, version, buildDate)
		}
	}
	app.Description = func() string { return "A set of utilities for processing audio files for retro console development" }
	app.Banner = func() {
		fmt.Fprintln(os.Stderr, "\n\033[1;97mPlease help support the Retcon project at https://www.undeveloper.com/retcon#support\033[0m")
	}
	app.RegisterCommand("psg", psg.Cmd)

	err := app.Run(os.Args[1:])
	if err != nil {
		os.Exit(1)
		return
	}
	os.Exit(0)
}
