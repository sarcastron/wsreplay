/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
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
		config, err := appConfig.GetConfig(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Target: %s\n", output.Info(config.Target))
		fmt.Printf("Duration: %s\n", output.Info(config.Duration))
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
