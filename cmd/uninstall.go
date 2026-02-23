/*
Copyright © 2026 Faissal Maulana
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "A command to uninstall package's config files",
	Long: `uninstall will deletes package's config files
`,
	Aliases: []string{"u"},
	RunE: func(cmd *cobra.Command, args []string) error {
		withExac, err := cmd.Flags().GetBool("ex")
		if err != nil {
			return err
		}

		if withExac {
			fmt.Println("uninstall called to delete configs of", args[0], "with exactly name package")
			return nil
		}

		fmt.Println("uninstall called to delete configs of", args[0])
		return nil
	},
}

func init() {
	uninstallCmd.Flags().BoolP("ex", "e", false, "typed exactly the package name")
	rootCmd.AddCommand(uninstallCmd)
}
