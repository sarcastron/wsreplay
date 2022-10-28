package wsrecorder

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func Record(uri string) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	startTime := time.Now()
	i := 0
	fmt.Printf("connecting to %s\n", uri)
	c, _, err := websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			i += 1
			fmt.Printf("%v T: %v - %s", i, time.Since(startTime), message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		// case t := <-ticker.C:
		// 	err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		// 	if err != nil {
		// 		log.Println("write:", err)
		// 		return
		// 	}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// type Recorder struct {
// 	client *websocket.Conn
// }

// func NewRecorder(uri *url.URL) *Recorder {
// 	conn, _, err := websocket.DefaultDialer.Dial(uri.String(), nil)
// 	if err != nil {
// 		log.Fatal("Dial: Error connecting to Websocket Server:", err)
// 	}
// 	return &Recorder{conn}
// }

// func (r Recorder) Close() {
// 	r.client.Close()
// }

// func (r Recorder) Record(seconds int) {
// 	//
// }
