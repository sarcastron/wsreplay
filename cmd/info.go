/*
Copyright © 2022 Adam Plante <toomanyadams@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	appConfig "wsreplay/pkg/config"
	"wsreplay/pkg/output"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get info about the app configuration.",
	Long:  `Displays info about the applications configuration if one is detected.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile == "" {
			fmt.Println("No configuration file provided.")
			os.Exit(1)
		}
		config, err := appConfig.LoadConfig(&cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Print("Recorder settings\n")
		fmt.Printf(" - Target: %s\n", output.Info(config.Target))
		fmt.Printf(" - Duration: %s\n", output.Info(config.Duration))
		fmt.Printf(" - Output File: %s\n", output.Info(config.File))

		fmt.Print("\nPlayback settings\n")
		fmt.Printf(" - Input File: %s\n", output.Info(config.File))
		fmt.Printf(" - Server Address: %s\n", output.Info(config.ServerAddr))
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
