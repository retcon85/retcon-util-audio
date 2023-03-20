# Retcon audio utilities

[![Build](https://github.com/retcon85/retcon-util-audio/actions/workflows/build.yml/badge.svg)](https://github.com/retcon85/retcon-util-audio/actions/workflows/build.yml)

## Overview

Part of the [Retcon85](https://github.com/retcon85) project, this is one of a set of tools aiming to provide a low barrier to entry to learning about computers by building a console:

- The [Retcon kit](https://www.undeveloper.com/retcon) itself and accompanying literature
- The [toolchain-sms](https://github.com/retcon85/toolchain-sms) Docker image providing a "one-click" toolchain for Sega Master System/Game Gear development
- The [template-sms-devkitsms](https://github.com/retcon85/template-sms-devkitsms) and [template-sms-wladx](https://github.com/retcon85/template-sms-wladx) template repositories for rapid boostrapping of C and Z80 assembler projects respectively
- This utility and related [retcon-util-gfx](https://github.com/retcon85/retcon-util-gfx) utilities for graphics handling

Currently this tool only supports the [PSG format](https://github.com/sverx/PSGlib/blob/master/documents/PSG%20file%20format.txt) for SMS/GG but as more canonical project builds are added it is designed to expand through new sub-commands.

The most useful things you can currently do with this tool is compress PSG files as well as visually debug them.

To convert from VGM to PSG format, it's recommended to use the `vgm2psg` tool under [PSGLib](https://github.com/sverx/PSGlib/tree/master/tools), which is also available on `toolchain-sms`

If you're looking for a tracker to compose VGM music / sound effects, why not take a look at [the Furnace project](https://github.com/tildearrow/furnace)

## Installation

### As part of toolchain-sms

The `toolchain-sms` docker image contains this utility from version 0.9 onwards.

```
docker pull retcon85/toolchain-sms
```

### Binaries

Download the latest releases from https://github.com/retcon85/retcon-util-audio/releases

### Building from source

[Go](https://go.dev/dl/) is the only prerequisite.

1. Clone the repo
2. Run `go . run` to run
3. Run `go . build` to build

Optionally, if you have GNU make installed you can run the Makefile.

## Usage

### Global usage

```
Usage:

      retcon-audio <command> [options]

Available commands:

      psg           utilities for processing SMS PSG files

Global options:

      --debug       prints extra debug information for selected commands
  -h, --help        prints help about a command
      --no-banner   suppresses the banner text after this program runs
      --quiet       suppresses all output except errors and banner
      --silent      suppresses all output except errors
```

### The `psg` sub-command usage

```
Usage:

  retcon-util-audio psg [options]
    see below for options reference

  retcon-util-audio psg [options] <psgfile>
    equivalent to "retcon-util-audio psg <options> --load <psgfile>"

  retcon-util-audio psg [options] <psgfile> <outfile>
    equivalent to "retcon-util-audio psg <options> --load <psgfile> --save <outfile>"

Options:

  -f, --load string     path of a PSG file to load, or "-" to read from standard input (default "-")
  -u, --no-compress     do not compress PSG output
  -o, --output string   output format. one of: psg | debug[=...] (default "psg")
  -s, --save string     path of the output file to generate, or "-" to write to standard output (default "-")

Debug format options:

  --output debug=(o|b|f|*)...
    o - print offset of byte from input source
    b - print raw byte(s) from input source
    f - print frame end markers
    * - print with all of the above options enabled

  --output debug
    print with none of the above options enabled

Global options:

      --debug       prints extra debug information for selected commands
  -h, --help        prints help about a command
      --no-banner   suppresses the banner text after this program runs
      --quiet       suppresses all output except errors and banner
      --silent      suppresses all output except errors
```
