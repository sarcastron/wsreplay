package wsserver

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

const pongWait = 10 * time.Second

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func StartServer(addr string, wsChan chan *websocket.Conn) {
	wsConnectionHandler := func(w http.ResponseWriter, r *http.Request) {
		// TODO only allow specific origins
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		ws.SetPongHandler(func(string) error {
			log.Println("Got the pong")
			ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		go func() {
			for {
				_, message, err := ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
					}
					break
				}
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				log.Printf("recv: %s", message)
			}
		}()
		wsChan <- ws
	}
	http.HandleFunc("/", wsConnectionHandler)
	// Wrap this in a go routine so it doesn't block.
	go func() {
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
}
