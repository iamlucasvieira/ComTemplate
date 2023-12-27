/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/iamlucasvieira/ComTemplate/pkg/cli"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a default template file",
	Long: `Creates a default template file named 'comtemplate.yml'
    at the current directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cli.CreateDefault()
		if err != nil {
			cli.Write(
				cli.Header("Error creating default file"),
				err.Error(),
			)
			os.Exit(1)
		}
		cli.Write(
			"Default file created",
		)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
