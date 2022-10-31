/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"wsreplay/pkg/tapedeck"
)

var playbackFile string

// playbackCmd represents the playback command
var playbackCmd = &cobra.Command{
	Use:   "playback",
	Short: "Playback a recorded websocket session.",
	Long:  `Will playback a recorded session. Playback will start as soon as the client connects to it unless the --immediate flag is set.`,
	Run: func(cmd *cobra.Command, args []string) {
		var messages []tapedeck.Message
		err := tapedeck.ReadTape(playbackFile, &messages)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tapedeck.Playback(&messages)
		// fmt.Println(messages)
	},
}

func init() {
	rootCmd.AddCommand(playbackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playbackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playbackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	playbackCmd.Flags().StringVarP(&playbackFile, "file", "f", "", "The file to playback.")
}
