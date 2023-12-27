/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	"github.com/iamlucasvieira/ComTemplate/pkg/cli"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ct",
	Short: "ComTemplate (ct) makes it easy to create and use git commit templates",
	Long: `With ComTemplate you can create and use git commit templates.

1. Initialize a default template file: 'ct init'
- This will create a file named 'comtemplate.yml' at the current directory 
containing default templates of git commit messages. You can edit this file to
add your own templates.

2. List available templates: 'ct list'
- This will list all available templates in the 'comtemplate.yml' file.

3. Use a template: 'ct <template-name>'
- This will open a form to fill the template variables. After filling the form,
the commit message will be printed to the terminal and copied to the clipboard.
You can paste it in your commit message.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data := getTemplates()
		t, ok := data[args[0]]
		if !ok {
			fmt.Printf("Template '%s' not found\n", args[0])
			os.Exit(1)
		}
		headerStr := fmt.Sprintf("Using template '%s'", t.Name)
		cli.Write(
			cli.Header(headerStr),
		)
		text, err := cli.PopulateFromForm(t)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = clipboard.WriteAll(text)
		if err != nil {
			cli.Write(
				cli.Header("Error copying to clipboard"),
				err.Error(),
			)
			os.Exit(1)
		}

		cli.WriteNoMargin(
			text,
			cli.TextHighlight("✔ Copied to clipboard"),
		)
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

func getTemplates() map[string]cli.Template {
	data, err := cli.ReadDefault()
	if err != nil {
		fmt.Println(`Error reading default file

Make sure you have a file named 'comtemplate.yml' or 'comtemplate.yaml' at the current directory.

Run: 'comtemplate init' to create a default file.
        `)
		os.Exit(1)
	}

	// Turn into map
	t := make(map[string]cli.Template)
	for _, template := range data {
		t[template.Name] = template
	}
	return t
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
