package wsserver

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func makeWsConnectionHandler(wsChan chan *websocket.Conn) func(w http.ResponseWriter, r *http.Request) {
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
		wsChan <- ws
	}
	return wsConnectionHandler
}

func StartServer(addr string, wsChan chan *websocket.Conn) {
	http.HandleFunc("/", makeWsConnectionHandler(wsChan))
	// Wrap this in a go routine so it doesn't block.
	go func() {
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
}
