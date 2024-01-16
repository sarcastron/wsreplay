package tapedeck

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"wsreplay/pkg/output"

	"github.com/gorilla/websocket"
)

func PlaybackAsync(messages *[]Message, wsConn *websocket.Conn) chan BusMessager {
	msgBus := make(chan BusMessager)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		defer close(msgBus)
		startTime := time.Now()
		total := len(*messages)
		i := 0
		output.MeterSpinner.Spin()
	playbackLoop:
		for {
			select {
			case <-time.After(time.Millisecond / 2):
				ts := time.Since(startTime)
				msgBus <- &PlaybackPrompt{output.MeterSpinner.Render(), ts}
				if ts >= (*messages)[i].Tick {
					msgBus <- &BusMessageInfo{"", fmt.Sprintf("#%d - %s", i+1, strings.TrimSuffix(string((*messages)[i].Content), "\n"))}
					// TODO Allow sending binary ws messages as well.
					err := wsConn.WriteMessage(websocket.TextMessage, (*messages)[i].Content)
					if err != nil {
						msgBus <- &BusMessageErr{"Tx :", err, false}
						fmt.Println(" Connection error. Shutting down...")
						break playbackLoop
					}
					i += 1
					if i >= total {
						break playbackLoop
					}
				}
			case <-interrupt:
				fmt.Println(" Interrupt signal detected. Shutting down...")
				break playbackLoop
			}
		}
	}()
	return msgBus
}
