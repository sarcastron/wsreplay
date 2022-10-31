package tapedeck

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"wsreplay/pkg/output"

	"github.com/gorilla/websocket"
)

func Playback(messages *[]Message, wsConn *websocket.Conn) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})
	defer close(done)

	// ticker := time.NewTicker(time.Millisecond / 2)

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

	for {
		select {
		case <-done:
			fmt.Println(" - unexpected Done.")
			return nil
		case <-time.After(time.Millisecond / 2):
			// Check for duration to expire
			ts := time.Since(startTime)
			fmt.Printf("%v | %v - %d\r", (*messages)[i].Tick, ts, i)
			if ts >= (*messages)[i].Tick {
				fmt.Println(" -- ", ts, (*messages)[i].Tick, "--")
				// TODO Allow recording binary ws messages as well.
				err := wsConn.WriteMessage(websocket.TextMessage, (*messages)[i].Content)
				if err != nil {
					log.Println("write:", output.Danger(err))
				}
				i += 1
				if i >= total {
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
