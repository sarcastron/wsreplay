/*
Copyright © 2022 Adam Plante <toomanyadams@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"

	appConfig "wsreplay/pkg/config"
	"wsreplay/pkg/tapedeck"
	"wsreplay/pkg/wsserver"
)

var playbackFile string

// playbackCmd represents the playback command
var playbackCmd = &cobra.Command{
	Use:   "playback",
	Short: "Playback a recorded websocket session.",
	Long:  `Will playback a recorded session. Playback will start as soon as the client connects to it unless the --immediate flag is set.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := appConfig.GetConfig(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var messages []tapedeck.Message
		err = tapedeck.ReadTape(playbackFile, &messages)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(messages) < 1 {
			fmt.Println(" -- No messages to playback.")
			os.Exit(0)
		}
		// This channel will be used to pass back the websocket connection ref
		wsChan := make(chan *websocket.Conn)
		// Start the http server
		wsserver.StartServer(config.ServerAddr, wsChan)
		fmt.Println("Server is listening on ", config.ServerAddr)
		fmt.Print("Waiting for a client to connect...")
		wsConn := <-wsChan
		defer wsConn.Close()
		fmt.Println("connected.")
		tapedeck.Playback(&messages, wsConn)
	},
}

func init() {
	RootCmd.AddCommand(playbackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playbackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playbackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	playbackCmd.Flags().StringVarP(&playbackFile, "file", "f", "", "The file to playback.")
}
