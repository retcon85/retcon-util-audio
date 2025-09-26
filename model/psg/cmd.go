package psg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/willbritton/gocli"
)

var Cmd = gocli.Command{
	Description: "utilities for processing SMS PSG files",
	Handler: func(cli *gocli.Cli, cmdArg string, arguments []string) error {
		fullCmd := fmt.Sprintf("%s %s", cli.Name, cmdArg)
		args := flag.NewFlagSet(fullCmd, flag.ContinueOnError)
		args.Usage = func() {
			// fmt.Fprintf(os.Stderr, "%s: %s\n\n", fullCmd, c.Description())
			fmt.Fprintln(os.Stderr, "Usage:")
			fmt.Fprintln(os.Stderr)
			fmt.Fprintf(os.Stderr, "  %[1]s %[2]s [options]\n", cli.Name, cmdArg)
			fmt.Fprintln(os.Stderr, "    see below for options reference")
			fmt.Fprintf(os.Stderr, "\n  %[1]s %[2]s [options] <psgfile>\n", cli.Name, cmdArg)
			fmt.Fprintf(os.Stderr, "    equivalent to \"%s %s <options> --load <psgfile>\"\n\n", cli.Name, cmdArg)
			fmt.Fprintf(os.Stderr, "  %[1]s %[2]s [options] <psgfile> <outfile>\n", cli.Name, cmdArg)
			fmt.Fprintf(os.Stderr, "    equivalent to \"%s %s <options> --load <psgfile> --save <outfile>\"\n\n", cli.Name, cmdArg)
			fmt.Fprintln(os.Stderr, "Options:")
			fmt.Fprintln(os.Stderr)
		}
		load := args.StringP("load", "f", "-", "path of a PSG file to load, or \"-\" to read from standard input")
		save := args.StringP("save", "s", "-", "path of the output file to generate, or \"-\" to write to standard output")
		format := args.StringP("output", "o", "psg", "output format. one of: psg | debug[=...]")
		noCompress := args.BoolP("no-compress", "u", false, "do not compress PSG output")

		cli.IgnoreGlobalOptions(args, []string{"help"})

		err := args.Parse(arguments)

		if err == flag.ErrHelp {
			args.PrintDefaults()
			printDebugOptions()
			cli.PrintGlobalOptions()
		}
		if err != nil {
			return err
		}

		posArgs := args.Args()
		if len(posArgs) > 2 {
			return handleError(args, errors.New("too many positional arguments"))
		}
		if len(posArgs) >= 1 {
			*load = posArgs[0]
		}
		if len(posArgs) == 2 {
			*save = posArgs[1]
		}

		// buffer the input in case load and save point to the same file
		src := new(bytes.Buffer)
		src.ReadFrom(getInReader(*load))
		dst := getOutWriter(*save)

		switch {
		case *format == "psg":
			if *noCompress {
				Decompress(src, dst)
			} else {
				buf := new(bytes.Buffer)
				Decompress(src, buf)
				Compress(buf, dst)
			}
		case (*format)[:5] == "debug":
			dbgFmt := (*format)[5:]
			Debug(src, dst, DebugOptions{
				PrintOffset: strings.ContainsAny(dbgFmt, "o*"),
				PrintBytes:  strings.ContainsAny(dbgFmt, "b*"),
				ShowFrames:  strings.ContainsAny(dbgFmt, "f*"),
			})
		default:
			return handleError(args, fmt.Errorf("unrecognised output format '%s'", *format))
		}

		return nil
	},
}

func printDebugOptions() {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, `Debug format options:

  --output debug=(o|b|f|*)...
    o - print offset of byte from input source
    b - print raw byte(s) from input source
    f - print frame end markers
    * - print with all of the above options enabled

  --output debug
    print with none of the above options enabled`)
}

func handleError(args *flag.FlagSet, err error) error {
	fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
	args.Usage()
	args.PrintDefaults()
	printDebugOptions()
	return err
}

func getInReader(flag string) io.Reader {
	if flag == "-" {
		return os.Stdin
	}
	f, err := os.Open(flag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file '%s' for reading\n", flag)
		os.Exit(1)
	}
	return f
}

func getOutWriter(flag string) io.Writer {
	if flag == "-" {
		return os.Stdout
	}
	f, err := os.Create(flag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file '%s' for writing\n", flag)
		os.Exit(1)
	}
	return f
}
