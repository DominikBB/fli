/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/dominikbb/fli/internal/app"
	"github.com/dominikbb/fli/internal/history"
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
		}

		if err := a.Creator.Store(last, args...); err != nil {
			fmt.Println("failed to store the command", err.Error())
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
