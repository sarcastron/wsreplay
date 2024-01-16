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
	"wsreplay/pkg/output"
	"wsreplay/pkg/tapedeck"
	"wsreplay/pkg/wsserver"
)

var (
	file       *string
	serverAddr *string
)

// playbackCmd represents the playback command
var playbackCmd = &cobra.Command{
	Use:   "playback",
	Short: "Playback a recorded websocket session.",
	Long:  `Will playback a recorded session. Playback will start as soon as the client connects to it.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("quiet") {
			quiet = true
		}
		config, err := appConfig.LoadConfig(&cfgFile)
		if err != nil {
			output.ErrorMsg(err)
			os.Exit(1)
		}

		if config == nil {
			config, err = appConfig.NewPlaybackConfig(
				file,
				serverAddr,
			)
			if err != nil {
				output.ErrorMsg(err)
				os.Exit(1)
			}
		}

		fmt.Printf("Loading %s\n", output.Notice(config.File))

		var messages []tapedeck.Message
		err = tapedeck.ReadTape(config.File, &messages)
		if err != nil {
			output.ErrorMsg(err)
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
		msgBus := tapedeck.PlaybackAsync(&messages, wsConn)
		for msg := range msgBus {
			switch bm := msg.(type) {
			case *tapedeck.PlaybackPrompt:
				fmt.Print(msg.CliMessage())
			case *tapedeck.BusMessageInfo:
				if !quiet {
					fmt.Print(msg.CliMessage() + "                 \n")
				}
			case *tapedeck.BusMessageErr:
				if bm.IsFatal {
					fmt.Println(bm.CliMessageln())
					os.Exit(1)
				} else {
					fmt.Println(bm.CliMessageln())
				}
			}
		}
		fmt.Println("------------------------------------------")
		fmt.Printf("%s Messages replayed.\n", output.Info(len(messages)))
	},
}

func init() {
	RootCmd.AddCommand(playbackCmd)

	file = playbackCmd.Flags().StringP("file", "f", "", "The file to playback.")
	serverAddr = playbackCmd.Flags().StringP("server", "s", ":8001", "The address for the server.")
}
