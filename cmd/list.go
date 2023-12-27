/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/iamlucasvieira/ComTemplate/pkg/cli"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all templates",
	Long: `Returns a list with the name of each template found in
    the template configuration file.

    File is either 'comtemplate.yml' or 'comtemplate.yaml' at the current
    directory.
    `,
	Run: func(cmd *cobra.Command, args []string) {
		data := getTemplates()
		titles := []string{}
		for _, template := range data {
			titleAndDescription := fmt.Sprintf("%s: %s", template.Name, template.Description)
			titles = append(titles, titleAndDescription)
		}

		fmt.Println(cli.RenderList("Available templates", titles))

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
