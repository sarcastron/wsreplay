package tapedeck

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"wsreplay/pkg/output"

	"github.com/gorilla/websocket"
)

func PlaybackSync(messages *[]Message, wsConn *websocket.Conn) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})
	defer close(done)

	gracefulShutdown := func() {
		// Cleanly close everything down
		select {
		case <-done:
			println("Expected done.")
		case <-time.After(time.Second):
		}
	}

	startTime := time.Now()
	total := len(*messages)
	i := 0

	output.MeterSpinner.Spin()

	for {
		select {
		case <-done:
			fmt.Println(" - unexpected Done.")
			return nil
		case <-time.After(time.Millisecond / 2):
			// Check for duration to expire
			ts := time.Since(startTime)
			fmt.Printf("  %s %s          \r", output.Info(output.MeterSpinner.Render()), ts)
			if ts >= (*messages)[i].Tick {
				fmt.Printf("#%d - %s                 \n", i+1, strings.TrimSuffix(string((*messages)[i].Content), "\n"))
				// TODO Allow recording binary ws messages as well.
				err := wsConn.WriteMessage(websocket.TextMessage, (*messages)[i].Content)
				if err != nil {
					log.Println("write:", output.Danger(err))
				}
				i += 1
				if i >= total {
					fmt.Println("------------------------------------------")
					fmt.Printf("%s Messages replayed.\n", output.Info((total)))
					gracefulShutdown()
					return nil
				}
			}
		case <-interrupt:
			fmt.Println(" Interrupt signal detected. Shutting down...")
			gracefulShutdown()
			return nil
		}
	}
}