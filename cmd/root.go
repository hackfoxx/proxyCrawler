package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)
import "github.com/spf13/cobra"

var (
	configFlag string
)

var rootCmd = &cobra.Command{
	Use: filepath.Base(os.Args[0]),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(-1)
	},
}
var logo = `` +
	"  ___                   ___                 _         \n" +
	" | _ \\_ _ _____ ___  _ / __|_ _ __ ___ __ _| |___ _ _ \n" +
	" |  _/ '_/ _ \\ \\ / || | (__| '_/ _` \\ V  V / / -_) '_|\n" +
	" |_| |_| \\___/_\\_\\\\_, |\\___|_| \\__,_|\\_/\\_/|_\\___|_|\n" + `
[ProxyCrawler]                   By:hackfoxx   v1.0.0
======================================================`

/*const logoTemplate = `
{{.UsageString}}`*/

func init() {
	//rootCmd.SetHelpTemplate(logoTemplate)
	startCmd.Flags().StringVarP(&configFlag, "config", "c", "./config.yml", "Path to config file")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(validCmd)
	rootCmd.AddCommand(checkCmd)
	cobra.OnInitialize(beforeStart)
}

func beforeStart() {
	fmt.Println(logo)
}
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("cobra error: ", err.Error())
		os.Exit(1)
	}
}
