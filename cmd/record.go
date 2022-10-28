/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	appConfig "wsreplay/pkg/config"
	"wsreplay/pkg/output"
	"wsreplay/pkg/wsrecorder"

	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("record called")
		config, err := appConfig.GetConfig(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Recording: %s for %s seconds.\n", output.Info(config.Target), output.Info(config.Duration))
		wsrecorder.Record(config.Target)
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
