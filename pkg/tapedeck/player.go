package tapedeck

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"wsreplay/pkg/output"
)

func Playback(messages *[]Message) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})
	defer close(done)

	// ticker := time.NewTicker(time.Millisecond / 2)

	gracefulShutdown := func() {
		// Cleanly close everything down
		select {
		case <-done:
			println("Expected done.")
		case <-time.After(time.Second):
		}
	}

	startTime := time.Now()
	total := len(*messages)
	i := 0

	for {
		select {
		case <-done:
			fmt.Println(" - unexpected Done.")
			return nil
		case <-time.After(time.Millisecond / 2):
			// Check for duration to expire
			ts := time.Since(startTime)
			fmt.Printf("%v | %v - %d\r", (*messages)[i].Tick, ts, i)
			if ts >= (*messages)[i].Tick {
				fmt.Println(" -- ", ts, (*messages)[i].Tick, "--")
				i += 1
				if i >= total {
					fmt.Printf("%s Messages replayed.\n", output.Info((total)))
					gracefulShutdown()
					return nil
				}
			}
		case <-interrupt:
			fmt.Println(" Interrupt signal detected. Shutting down...")
			gracefulShutdown()
			return nil
		}
	}
}
