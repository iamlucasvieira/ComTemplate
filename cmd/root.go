/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.design/x/clipboard"

	"github.com/iamlucasvieira/ComTemplate/pkg/cli"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ComTemplate",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data := getTemplates()
		t := getTemplate(args[0], &data)

		if t.Name == "" {
			fmt.Printf("Template '%s' not found\n", args[0])
			os.Exit(1)
		}

		text, err := cli.PopulateFromForm(t)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = clipboard.Init()
		if err != nil {
			fmt.Println("Error copying to clipboard")
			fmt.Println(text)
			os.Exit(1)
		}

		clipboard.Write(clipboard.FmtText, []byte(text))

		fmt.Println(text)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// GetTemplates returns a list of templates
func getTemplates() []cli.Template {
	data, err := cli.ReadDefault()
	if err != nil {
		fmt.Println(`Error reading default file

Make sure you have a file named 'comtemplate.yml' or 'comtemplate.yaml' at the current directory.

Run: 'comtemplate init' to create a default file.
        `)
		os.Exit(1)
	}

	// Turn into map
	return data
}

// GetTemplate returns a template given its name
func getTemplate(name string, data *[]cli.Template) cli.Template {
	for _, t := range *data {
		if t.Name == name {
			return t
		}
	}

	return cli.Template{}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ComTemplate.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
