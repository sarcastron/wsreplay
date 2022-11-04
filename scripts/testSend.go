package main

import (
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type TradeMessage struct {
	Ticker string  `json:"t"`
	Price  float32 `json:"p"`
}

func RandPriceMovement() float32 {
	rand.Seed(time.Now().UnixNano())
	var min float32 = -2.0
	var max float32 = 2.0
	return min + rand.Float32()*(max-min)
}

func SendMessages() {
	ticker := "AAPL"
	u := url.URL{Scheme: "ws", Host: "localhost:8001", Path: "/", RawQuery: ""}
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
		r := rand.Intn(250)
		time.Sleep(time.Duration(r) * time.Millisecond)
		// time.Sleep(100 * time.Millisecond)

		message := TradeMessage{ticker, newPrice}
		if err := c.WriteJSON(message); err != nil {
			log.Fatalf("Failed to write JSON: %v\n", err)
		}
	}
}
