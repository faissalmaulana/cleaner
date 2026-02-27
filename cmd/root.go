/*
Copyright © 2026 Faissal Maulana
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var Version = "v1.2.3"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cleaner",
	Short: "A Command Line program to cleanup xdg configs",
	Long: `cleaner is a command line program to cleanup xdg configs a package after uninstallation of the package.

It will find common config files in $HOME and
delete all associate config files of the package.
	`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}
