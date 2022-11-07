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

func RecordAsync(uri string, duration time.Duration, messages *[]Message) chan BusMessager {
	msgBus := make(chan BusMessager)

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
				msgBus <- &BusMessageErr{"Rx:", err, false}
				return
			}
			ts := time.Since(startTime)
			*messages = append(*messages, Message{ts, message})
			i += 1
			msgBus <- &BusMessageInfo{fmt.Sprint(i), fmt.Sprintf("T: %v - %s", time.Since(startTime), message)}
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
			msgBus <- &BusMessageInfo{"", "Closing after time out."}
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
					// msgBus <- newError(err)
					msgBus <- &BusMessageErr{"Tx:", err, false}
				}
			case t := <-ticker.C:
				if duration != 0 && t.After(endTime) {
					// msgBus <- newMessage(fmt.Sprintf("\nDuration of %v has elapsed. Shutting down...\n", output.Notice(duration)))
					msgBus <- &BusMessageInfo{"", fmt.Sprintf("\nDuration of %s has elapsed. Shutting down...\n", duration)}
					gracefulShutdown()
					break recordLoop
				}
			case <-interrupt:
				msgBus <- &BusMessageInfo{"", "Interrupt signal (ctrl-c) detected. Shutting down..."}
				gracefulShutdown()
				break recordLoop
			}
		}
	}()

	return msgBus
}
