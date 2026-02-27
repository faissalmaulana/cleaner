/*
Copyright © 2026 Faissal Maulana
*/
package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	fp "github.com/faissalmaulana/cleaner/internal/filepath"
	"github.com/faissalmaulana/cleaner/internal/utils"
	"github.com/spf13/cobra"
)

// xdgDirs are the standard XDG user directories to search for package configs.
var XDGDirs = []string{".config", ".cache"}

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "A command to uninstall package's config files",
	Long: `uninstall will deletes package's config files
`,
	Aliases: []string{"u"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.PrintErrf("package name is required\n")
			return
		}

		err := utils.ValidatePkgName(args[0])
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			return
		}

		homedir, err := os.UserHomeDir()
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			return
		}

		withExac, err := cmd.Flags().GetBool("ex")
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			return
		}

		var getFilepaths = func(path, pkgName string) (string, error) {
			return fp.GetFilePathFromOS(path, pkgName, false)
		}

		if withExac {
			getFilepaths = func(path, pkgName string) (string, error) {
				return fp.GetFilePathFromOS(path, pkgName, true)
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh,
			syscall.SIGINT,  // Ctrl+C
			syscall.SIGTERM, // kill / systemd stop
			syscall.SIGHUP,  // terminal closed
		)

		// goroutine watches for signals
		go func() {
			select {
			case _ = <-sigCh:
				fmt.Printf("\nCancelling...\n")
				cancel()
				defer os.Exit(1)
			case <-ctx.Done():
				// just exit it
			}
			signal.Stop(sigCh)
		}()

		pkgfilepaths, err := fp.GetFilePaths(ctx, getFilepaths, homedir, args[0], XDGDirs)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				fmt.Println("operation cancelled")
				return
			}
			cmd.PrintErrf("%v\n", err)
			return
		}

		fmt.Print("Do you want to continue? (y/N): ")

		reader := bufio.NewReader(os.Stdin)

		var confirmed bool

		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				cmd.PrintErrf("%v\n", err)
				return
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
			return
		}

		err = fp.DeleteFilePaths(ctx, os.Remove, pkgfilepaths)
		if err != nil {
			cmd.PrintErrf("%v\n", err)
			return
		}

		fmt.Printf("Done deleted configs of %s\n", args[0])
	},
}

func init() {
	uninstallCmd.Flags().BoolP("ex", "e", false, "typed exactly the package name")
	rootCmd.AddCommand(uninstallCmd)
}
