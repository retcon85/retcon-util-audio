package main

import (
	"fmt"
	"io"
	"os"
	"path"

	flag "github.com/spf13/pflag"

	"github.com/retcon85/retcon-util-audio/logger"
	"github.com/retcon85/retcon-util-audio/model/psg"
)

func getInReader(flag string, arg string) io.Reader {
	var src io.Reader
	if flag == "" {
		flag = arg
	}
	switch flag {
	case "":
		fmt.Fprintln(os.Stderr, "must specify an input file")
		os.Exit(1)
	case "-":
		src = os.Stdin
	default:
		f, err := os.Open(flag)
		src = f
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening file '%s' for reading\n", flag)
			os.Exit(1)
		}
	}
	return src
}

func getOutWriter(flag string, arg string) io.Writer {
	var dst io.Writer
	if flag == "" {
		flag = arg
	}
	switch flag {
	case "":
		fmt.Fprintln(os.Stderr, "must specify an output file")
		os.Exit(1)
	case "-":
		dst = os.Stdout
	default:
		f, err := os.Create(flag)
		dst = f
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening file '%s' for writing\n", flag)
			os.Exit(1)
		}
	}
	return dst
}

func handleVerbose(verbose bool) {
	if verbose {
		l := logger.DefaultLogger()
		l.SetVerbose()
	}
}

func printUsage(cmd string) {
	fmt.Fprintf(os.Stderr, "Usage of %s %s:\n", prog, cmd)
}

var prog string = path.Base(os.Args[0])

func printCommandList() {
	fmt.Fprintf(os.Stderr, "  %8s psg compress\tcompresses a PSG file\n", prog)
	fmt.Fprintf(os.Stderr, "  %8s psg decompress\tdecompresses a PSG file\n", prog)
	fmt.Fprintf(os.Stderr, "  %8s psg debug\tprints debug information about a PSG file\n", prog)
	fmt.Fprintf(os.Stderr, "\nRun %s (command) --help for more information\n", prog)
}

func main() {
	if len(os.Args) < 3 {
		printUsage("")
		printCommandList()
		os.Exit(1)
	}
	cmd := fmt.Sprintf("%s %s", os.Args[1], os.Args[2])
	switch cmd {
	case "psg compress":
		args := flag.NewFlagSet(fmt.Sprintf("%s %s", prog, cmd), flag.ContinueOnError)
		verbose := args.BoolP("verbose", "v", false, "logs debug information")
		in := args.StringP("input", "i", "", "the PSG file to compress, or \"-\" to read from standard input")
		out := args.StringP("output", "o", "", "the path of the compressed PSG file to generate, or \"-\" to write to standard output")
		err := args.Parse(os.Args[3:])
		if err == flag.ErrHelp {
			os.Exit(1)
		}
		handleVerbose(*verbose)
		src := getInReader(*in, args.Arg(0))
		dst := getOutWriter((*out), args.Arg(1))
		psg.Compress(src, dst)
	case "psg decompress":
		args := flag.NewFlagSet(fmt.Sprintf("%s %s", prog, cmd), flag.ContinueOnError)
		in := args.StringP("input", "i", "", "the PSG file to decompress, or \"-\" to read from standard input")
		out := args.StringP("output", "o", "", "the path of the decompressed PSG file to generate, or \"-\" to write to standard output")
		err := args.Parse(os.Args[3:])
		if err == flag.ErrHelp {
			os.Exit(1)
		}
		src := getInReader(*in, args.Arg(0))
		dst := getOutWriter((*out), args.Arg(1))
		psg.Decompress(src, dst)
	case "psg debug":
		args := flag.NewFlagSet(fmt.Sprintf("%s %s", prog, cmd), flag.ContinueOnError)
		in := args.StringP("input", "i", "", "the PSG file to analyze, or \"-\" to read from standard input")
		printOffset := args.BoolP("print-offset", "a", true, "prints the start address of each decoded line")
		printBytes := args.BoolP("print-bytes", "b", true, "prints the byte data for each decoded line")
		err := args.Parse(os.Args[3:])
		if err == flag.ErrHelp {
			os.Exit(1)
		}
		src := getInReader(*in, args.Arg(0))
		psg.Debug(src, os.Stdout, psg.DebugOptions{PrintOffset: *printOffset, PrintBytes: *printBytes})
	default:
		args := flag.NewFlagSet(prog, flag.ContinueOnError)
		err := args.Parse(os.Args[1:])
		if err == flag.ErrHelp {
			printCommandList()
			os.Exit(1)
			return
		}
		os.Exit(1)
	}
	os.Exit(0)
}
