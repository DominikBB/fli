/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dominikbb/fli/internal/app"
	"github.com/dominikbb/fli/internal/runner"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a command stored under a set of tags",
	Long:  `Run is used to run a particular command using a set of tags to identify it.`,
	Run: func(cmd *cobra.Command, args []string) {
		a := app.NewApp()

		if len(args) == 0 {
			ls, err := a.Getter.List()
			if err != nil {
				fmt.Printf("Unable to list commands %v\n", err)
				return
			}
			output, err := tea.NewProgram(runner.NewModel(runner.CommandsToItems(ls)), tea.WithAltScreen()).Run()
			if err != nil {
				fmt.Println("Error running program:", err)
				os.Exit(1)
			}
			if m, ok := output.(runner.Model); ok && m.SelectedCommand != "" {
				if err := runner.Run(m.SelectedCommand); err != nil {
					fmt.Printf("Failed to execute command: %v\n", err)
					return
				}
			}
			return
		}

		c, err := a.Getter.Get(args...)
		if err != nil {
			fmt.Printf("Unable to find command with tags %s: %v\n", args, err)
			return
		}
		if c == "" {
			fmt.Printf("Bad command at tags %s: %v\n", args, err)
			return
		}

		if err := runner.Run(c); err != nil {
			fmt.Printf("Failed to execute command: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
