/*
Copyright © 2022 Adam Plante <toomanyadams@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	appConfig "wsreplay/pkg/config"
	"wsreplay/pkg/output"
	"wsreplay/pkg/tapedeck"

	"github.com/spf13/cobra"
)

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record a websocket session.",
	Long:  `Records a websocket session and saves the session to serialized gob files.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := appConfig.GetConfig(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Recording: %s for %s seconds.\n", output.Info(config.Target), output.Info(config.Duration))
		var messages []tapedeck.Message
		tapedeck.Record(config.Target, time.Duration(config.Duration)*time.Second, &messages)
		fmt.Printf("%d message(s) recorded.\n", len(messages))
		err = tapedeck.WriteTape(config.OutputTapeFile, &messages)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(recordCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
