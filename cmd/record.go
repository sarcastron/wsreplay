/*
Copyright © 2022 Adam Plante <toomanyadams@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	appConfig "wsreplay/pkg/config"
	"wsreplay/pkg/output"
	"wsreplay/pkg/tapedeck"

	"github.com/spf13/cobra"
)

var target *string
var duration *int
var outputFile *string

// recordCmd represents the record command
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record a websocket session.",
	Long:  `Records a websocket session and saves the session to a serialized gob file.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := appConfig.LoadConfig(&cfgFile)
		if err != nil {
			output.ErrorMsg(err)
			os.Exit(1)
		}

		if config == nil {
			config, err = appConfig.NewRecordConfig(
				target,
				*duration,
				outputFile,
			)
			if err != nil {
				output.ErrorMsg(err)
				os.Exit(1)
			}
		}

		timeSpan := ""
		if *duration > 0 {
			timeSpan = fmt.Sprintf(" for %s seconds or", output.Info(config.Duration))
		}
		fmt.Printf("Recording: %s%s until interrupt (%s)\n", output.Info(config.Target), timeSpan, output.Notice("ctrl-c"))
		var messages []tapedeck.Message
		msgBus := tapedeck.RecordAsync(config.Target, time.Duration(config.Duration)*time.Second, &messages)
		for msg := range msgBus {
			fmt.Println("Got a message", msg)
			switch bm := msg.(type) {
			case *tapedeck.BusMessageInfo:
				fmt.Print(msg.CliMessage())
			case *tapedeck.BusMessageErr:
				if bm.IsFatal {
					fmt.Println(bm.CliMessage())
					os.Exit(1)
				} else {
					fmt.Println(bm.CliMessage())
				}
			}
		}
		fmt.Printf("%s message(s) recorded.\n", output.Info(len(messages)))
		err = tapedeck.WriteTape(config.File, &messages)
		if err != nil {
			output.ErrorMsg(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(recordCmd)

	target = recordCmd.Flags().StringP("target", "t", "", "Websocket connection to record.")
	duration = recordCmd.Flags().IntP("duration", "d", 0, "Number of seconds to record. 0 seconds will run until interrupted (ctrl-c).")
	outputFile = recordCmd.Flags().StringP("file", "f", "", "File to save the recorded data to.")
}
