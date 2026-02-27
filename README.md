# Cleaner

Cleaner is a command line tool for Linux XDG-based desktop to automaticly delete configs files of a package after uninstallation.

## Overview

This tool is designed to be used after uinstallation of a package, when all configs in `$HOME` want to be deleted. Cleaner will search the package's configs in `.config`, `.cache`, etc... and deletes all related config files of the package.

## Requirements

- Go v1.26 or later

## Installation

```bash
go install github.com/faissalmaulana/cleaner@latest
```

## Usage

```bash
cleaner unistall chrome
```

## Available Commands

```text
Usage:
  cleaner [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  uninstall   A command to uninstall package's config files

Flags:
  -h, --help      help for cleaner
  -v, --version   version for cleaner

Use "cleaner [command] --help" for more information about a command.
```

## Help Subcommands

You can run `--help` for each subcommands to know more details about the subcommand and their flags:

```bash
cleaner uninstall --help
```
