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

// TODO return an error and let the command handle the various error states
func Record(uri string, duration time.Duration, messages *[]Message) []Message {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	startTime := time.Now()
	endTime := startTime.Add(duration)
	i := 0

	c, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		log.Fatal(" - Dial Err:", output.Danger(err))
	}
	defer c.Close()

	done := make(chan struct{})

	// Naming this to signal intent
	readMessageLoop := func() {
		// If ReadMessage() errors, close done.
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println(" - Rx connection:", output.Danger(err))
				return
			}
			ts := time.Since(startTime)
			*messages = append(*messages, Message{ts, message})
			i += 1
			fmt.Printf("%v T: %v - %s", i, time.Since(startTime), message)
		}
	}
	// Throw this loop in a go routine to prevent blocking.
	go readMessageLoop()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	gracefulShutdown := func() {
		// Cleanly close the connection by sending a close message and then
		// waiting (with timeout) for the server to close the connection.
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println(" - Write close error:", err)
			return
		}
		select {
		case <-done:
			println("Expected done.")
		case <-time.After(time.Second):
		}
	}

	for {
		select {
		case <-done:
			fmt.Println(" - unexpected Done.")
			return *messages
		case t := <-ticker.C:
			// Check for duration to expire
			if duration != 0 && t.After(endTime) {
				fmt.Printf("\nDuration of %v has elapsed. Shutting down...\n", output.Notice(duration))
				gracefulShutdown()
				return *messages
			}
		case <-interrupt:
			fmt.Println(" Interrupt signal detected. Shutting down...")
			gracefulShutdown()
			return *messages
		}
	}
}
