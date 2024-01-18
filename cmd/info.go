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

func firstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

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
		if len(config.SendMessages) > 0 {
			fmt.Print(" - Send Messages:\n")
			for _, sm := range config.SendMessages {
				fmt.Printf("   - At %s seconds - %s...\n", output.Info(sm.At), output.Info(firstN(sm.Message, 20)))
			}
		} else {
			fmt.Printf(" - Send Messages: %s\n", output.Info("none"))
		}

		fmt.Print("\nPlayback settings\n")
		fmt.Printf(" - Input File: %s\n", output.Info(config.File))
		fmt.Printf(" - Server Address: %s\n", output.Info(config.ServerAddr))
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
