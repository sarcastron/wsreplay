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

func RecordAsync(uri string, duration time.Duration, messages *[]Message) (chan *string, chan bool, chan error) {
	msgBus := make(chan *string)
	rComplete := make(chan bool)
	errBus := make(chan error)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// This will be closed on websocket read error or close
	done := make(chan interface{})

	startTime := time.Now()
	endTime := startTime.Add(duration)
	i := 0

	c, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		log.Fatal(" - Dial Err:", output.Danger(err))
	}

	userInputChan := UserInput()
	go func() {
		// If ReadMessage() errors, close done.
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				// exits the loop if connection is killed from outside and inside.
				fmt.Println(" - Rx:", output.Danger(err))
				return
			}
			ts := time.Since(startTime)
			*messages = append(*messages, Message{ts, message})
			i += 1
			msg := fmt.Sprintf(" <- %v T: %v - %s", output.Notice(i), time.Since(startTime), message)
			msgBus <- &msg
		}
	}()

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
		case <-time.After(time.Second):
			msg := "Closing after timed out."
			msgBus <- &msg
		}
	}

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	go func() {
	recordLoop:
		for {
			select {
			case <-done:
				// Connection was closed unexpectedly. break the loop
				break recordLoop
			case input := <-userInputChan:
				fmt.Println("Sending message...")
				err := c.WriteMessage(websocket.TextMessage, []byte(*input))
				if err != nil {
					errBus <- err
				}
			case t := <-ticker.C:
				if duration != 0 && t.After(endTime) {
					msg := fmt.Sprintf("\nDuration of %v has elapsed. Shutting down...\n", output.Notice(duration))
					msgBus <- &msg
					gracefulShutdown()
					return
				}
			case <-interrupt:
				msg := " Interrupt signal detected. Shutting down...\n"
				msgBus <- &msg
				gracefulShutdown()
				break recordLoop
			}
		}
		rComplete <- true
		c.Close()
	}()

	return msgBus, rComplete, errBus
}
