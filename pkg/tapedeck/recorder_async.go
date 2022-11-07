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

type RecorderMessage struct {
	Content *string
	Err     error
}

func newMessage(msg string) *RecorderMessage {
	return &RecorderMessage{&msg, nil}
}
func newError(err error) *RecorderMessage {
	return &RecorderMessage{nil, err}
}

func RecordAsync(uri string, duration time.Duration, messages *[]Message) chan *RecorderMessage {
	msgBus := make(chan *RecorderMessage)

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
				msgBus <- newError(fmt.Errorf(" - Rx: %s", output.Danger(err)))
				return
			}
			ts := time.Since(startTime)
			*messages = append(*messages, Message{ts, message})
			i += 1
			msgBus <- newMessage(fmt.Sprintf(" <- %v T: %v - %s", output.Notice(i), time.Since(startTime), message))
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
			msgBus <- newMessage("Closing after time out.")
		}
	}

	ticker := time.NewTicker(time.Second)
	go func() {
		defer ticker.Stop()
		defer close(msgBus)
		defer c.Close()
	recordLoop:
		for {
			select {
			case <-done:
				// Connection was closed unexpectedly. break the loop
				break recordLoop
			case input := <-userInputChan:
				err := c.WriteMessage(websocket.TextMessage, []byte(*input))
				if err != nil {
					msgBus <- newError(err)
				}
			case t := <-ticker.C:
				if duration != 0 && t.After(endTime) {
					msgBus <- newMessage(fmt.Sprintf("\nDuration of %v has elapsed. Shutting down...\n", output.Notice(duration)))
					gracefulShutdown()
					break recordLoop
				}
			case <-interrupt:
				msgBus <- newMessage(" Interrupt signal detected. Shutting down...\n")
				gracefulShutdown()
				break recordLoop
			}
		}
	}()

	return msgBus
}
