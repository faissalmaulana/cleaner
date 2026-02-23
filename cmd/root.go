/*
Copyright © 2026 Faissal Maulana
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	VERSION string = "0.0.0"
)

var (
	version bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cleaner",
	Short: "A Command Line program to cleanup xdg configs",
	Long: `cleaner is a command line program to cleanup xdg configs a package after uninstallation of the package.

It will find common config files in $HOME and
delete all associate config files of the package.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if version {
			fmt.Println(VERSION)
			return nil
		}

		return nil
	},
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "installed cleaner version")

}
