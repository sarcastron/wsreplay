package tapedeck

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	"wsreplay/pkg/config"
	"wsreplay/pkg/output"

	"github.com/gorilla/websocket"
)

type SendableMessage struct {
	At      time.Time
	Message string
}

func RecordAsync(uri string, duration time.Duration, messages *[]Message, sendMessages []config.SendMessage) chan BusMessager {
	msgBus := make(chan BusMessager)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// This will be closed on websocket read error or when recordLoop is broken
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
			// Append to messages slice
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

	// Convert sendMessages to SendableMessage slice to be used for timed TX messages
	messagesWithTime := make([]SendableMessage, len(sendMessages))
	if len(sendMessages) > 0 {
		// messagesWithTime = make([]SendableMessage, len(sendMessages))
		for i, sm := range sendMessages {
			messagesWithTime[i] = SendableMessage{startTime.Add(time.Duration(sm.At * float32(time.Second))), sm.Message}
		}
	}
	var nextMessageIdx int = -1
	if len(messagesWithTime) > 0 {
		nextMessageIdx = 0
	}

	ticker := time.NewTicker(200 * time.Millisecond)
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
					msgBus <- &BusMessageErr{"Tx:", err, false}
				}
			case t := <-ticker.C:
				// check that there are messages to send and that the time has elapsed
				if nextMessageIdx != -1 && t.After(messagesWithTime[nextMessageIdx].At) {
					err := c.WriteMessage(websocket.TextMessage, []byte(messagesWithTime[nextMessageIdx].Message))
					if err != nil {
						msgBus <- &BusMessageErr{"Tx:", err, false}
					} else {
						msgBus <- &BusMessageInfo{"Tx", fmt.Sprintf("%s\n", messagesWithTime[nextMessageIdx].Message)}
					}
					// Set the next message index
					if nextMessageIdx < len(messagesWithTime)-1 {
						nextMessageIdx += 1
					} else {
						nextMessageIdx = -1
					}
				}
				if duration != 0 && t.After(endTime) {
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
