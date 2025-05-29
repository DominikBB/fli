/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
		}

		for _, command := range list {
			fmt.Printf("# %s -> %s\n", command[0], command[1])
		}
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
