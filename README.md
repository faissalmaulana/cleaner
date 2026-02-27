# Cleaner

Cleaner is a command line tool for Linux XDG-based desktop to automaticly delete configs files of a package after uninstallation.

## Overview

Cleaner helps you clean up leftover configuration files after uninstalling a package on Linux XDG-based desktops. It searches for and removes the package's config files in standard XDG directories like `.config`, `.cache` within `$HOME`.

## Requirements

- Go v1.26 or later

## Installation

```bash
go install github.com/faissalmaulana/cleaner@latest
```

## Usage

```bash
cleaner uninstall chrome
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
