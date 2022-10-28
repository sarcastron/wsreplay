package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type TradeMessage struct {
	Ticker string  `json:"t"`
	Price  float32 `json:"p"`
}

type TicketResponse struct {
	Ticket string `json:"ticket"`
}

func getTicket() string {
	api_key := os.Getenv("PUB_AUTH_KEY")
	url := "https://staging.api.prospero.ai/user/ticket"
	req, _ := http.NewRequest("Get", url, bytes.NewBuffer(make([]byte, 0)))
	req.Header.Add("Authorization", "Bearer "+api_key)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	var tr TicketResponse
	if err = json.Unmarshal(body, &tr); err != nil {
		log.Println("Error while unmarshaling response:", err)
	}
	return tr.Ticket
}

func RandPriceMovement() float32 {
	rand.Seed(time.Now().UnixNano())
	var min float32 = -2.0
	var max float32 = 2.0
	return min + rand.Float32()*(max-min)
}

func sendMessages(ticket string) {
	ticker := "AAPL"
	u := url.URL{Scheme: "ws", Host: "localhost:8001", Path: fmt.Sprintf("/v1/pub/price/%s", ticker), RawQuery: fmt.Sprintf("ticket=%s", ticket)}
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("WS Connection Failed: %v %v\n", err, resp)
	}
	defer c.Close()

	var startingPrice float32 = 123.45
	iterations := 2000
	log.Println("Going...")

	newPrice := startingPrice

	for x := 0; x < iterations; x++ {
		pm := RandPriceMovement()
		newPrice = newPrice + pm
		log.Printf("%d/%d %s -> %f (%f)\n", x+1, iterations, ticker, newPrice, pm)
		time.Sleep(100 * time.Millisecond)

		message := TradeMessage{ticker, newPrice}
		if err := c.WriteJSON(message); err != nil {
			log.Fatalf("Failed to write JSON: %v\n", err)
		}
	}
}

func main() {
	ticket := getTicket()
	log.Println(ticket)
	sendMessages(ticket)
}
