/*
Copyright © 2026 Faissal Maulana
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	fp "github.com/faissalmaulana/cleaner/internal/filepath"
	"github.com/faissalmaulana/cleaner/internal/utils"
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
		if len(args) == 0 {
			return errors.New("package name is required")
		}

		err := utils.ValidatePkgName(args[0])
		if err != nil {
			return err
		}

		homedir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		withExac, err := cmd.Flags().GetBool("ex")
		if err != nil {
			return err
		}

		var getFilepaths = func(path, pkgName string) (string, error) {
			return fp.GetFilePathFromOS(path, pkgName, false)
		}

		if withExac {
			getFilepaths = func(path, pkgName string) (string, error) {
				return fp.GetFilePathFromOS(path, pkgName, true)
			}
		}

		pkgfilepaths, err := fp.GetFilePaths(getFilepaths, homedir, args[0])
		if err != nil {
			return err
		}

		fmt.Print("Do you want to continue? (y/N): ")

		reader := bufio.NewReader(os.Stdin)

		var confirmed bool

		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			input = strings.TrimSpace(strings.ToLower(input))

			switch input {
			case "y", "yes":
				confirmed = true
			case "n", "no", "":
				confirmed = false
			default:
				fmt.Println("Please enter 'y' or 'n'.")
				continue
			}

			break
		}

		if !confirmed {
			fmt.Println("Aborted.")
			return nil
		}

		err = fp.DeleteFilePaths(os.Remove, pkgfilepaths)
		if err != nil {
			return err
		}

		fmt.Printf("Done deleted configs of %s\n", args[0])

		return nil
	},
}

func init() {
	uninstallCmd.Flags().BoolP("ex", "e", false, "typed exactly the package name")
	rootCmd.AddCommand(uninstallCmd)
}
