# cleaner is A Command Line Tool for linux xdg-based desktop to automaticly delete configs files of a package after uninstallation.

## this tool is designed to used after uinstallation of a package, and all configs that already in $HOME want to be deleted. Cleaner will search the package's configs in ```.config```, ```.cache```, etc ...  and deletes all related-config files of the package.


## Requirements:
 - Go v1.26 or later

## Installation: 
install cleaner by using:
```
go install github.com/faissalmaulana/cleaner@latest
```

## Usage:
```
cleaner uinstall chrome
```

## Available Commands:
```
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

## Help Subcommands:
you can run ```--help``` for each subcommands to know more detail about the subcommand and their flags
```
cleaner uninstall --help
```
