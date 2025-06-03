/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dominikbb/fli/internal/app"
	"github.com/dominikbb/fli/internal/history"
	"github.com/dominikbb/fli/internal/repository"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long:  `Push grabs the last command executed and pushes it into Fli history, under a set of specified tags.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		a := app.NewApp()
		h := history.History{Shell: history.ShellNu}
		last, err := h.GetLastCommand()
		if err != nil {
			fmt.Println("failed to find the last executed command")
			return
		}

		err = a.Creator.Store(last, args...)
		if err != nil && errors.Is(err, repository.ErrDuplicate) {
			existing, err := a.Getter.Get(args...)
			if err != nil {
				fmt.Println("failed to find the conflicting command")
				return
			}
			fmt.Println("A command is already assigned to these tags:")
			fmt.Printf("# %s -> %s\n", strings.Join(args, ", "), existing)
			fmt.Println("New command:")
			fmt.Printf("# %s -> %s\n", strings.Join(args, ", "), last)
			fmt.Print("Would you like to overwrite the existing command? (y/N): ")
			var response string
			_, scanErr := fmt.Scanln(&response)
			if scanErr != nil || (len(response) > 0 && (response[0] == 'y' || response[0] == 'Y')) {
				err := a.Remover.Delete(args...)
				if err != nil {
					fmt.Println("failed to overwrite the command", err.Error())
				}

				overwriteErr := a.Creator.Store(last, args...)
				if overwriteErr != nil {
					fmt.Println("failed to overwrite the command", overwriteErr.Error())
					return
				}

				fmt.Printf("Pushed command\n# %s -> %s", strings.Join(args, ", "), last)
				fmt.Printf("\n\nRepeat it by running:\n    fli run %s", strings.Join(args, " "))
			} else {
				fmt.Println("Did not overwrite the existing command.")
			}
			return
		}

		if err != nil {
			fmt.Println("failed to store the command", err.Error())
			return
		}

		fmt.Printf("Pushed command\n# %s -> %s", strings.Join(args, ", "), last)
		fmt.Printf("\n\nRepeat it by running:\n    fli run %s", strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// func getLastCommand(shell string) (string, error) {
// 	user, _ := user.Current()
// 	switch filepath.Base(shell) {
// 	case "bash", "sh":
// 		return readLastLine(filepath.Join(user.HomeDir, ".bash_history"))
// 	case "zsh":
// 		return readZshLastCommand(filepath.Join(user.HomeDir, ".zsh_history"))
// 	case "fish":
// 		return readFishLastCommand(filepath.Join(user.HomeDir, ".local/share/fish/fish_history"))
// 	case "nu":
// 		return readLastLine(filepath.Join(user.HomeDir, ".cache/nushell/history.txt"))
// 	default:
// 		return "", fmt.Errorf("unsupported shell")
// 	}
// }
