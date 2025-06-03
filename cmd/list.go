/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/dominikbb/fli/internal/app"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  `Lists all saved commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		a := app.NewApp()
		list, err := a.Getter.List(args...)
		if err != nil {
			fmt.Printf("Unable to list commands for tags %s: %v\n", args, err)
			return
		}

		if len(list) == 0 {
			fmt.Printf("No commands for tags: %s\n", args)
			return
		}

		var (
			purple = lipgloss.Color("99")
			gray   = lipgloss.Color("245")
			// lightGray = lipgloss.Color("241")

			headerStyle = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
			cellStyle   = lipgloss.NewStyle().Padding(0, 1)
		)

		t := table.New().
			Border(lipgloss.HiddenBorder()).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch row {
				case table.HeaderRow:
					return headerStyle
				default:
					if col == 0 {
						// First column: purple
						return cellStyle.Foreground(purple)
					}
					// Second column: gray
					return cellStyle.Foreground(gray)
				}
			}).
			Rows(list...)

		fmt.Println(t)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
