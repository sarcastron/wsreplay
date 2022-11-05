package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func prompt() {
	fmt.Print("-> ")
}

func InputGetter() chan *string {
	inputChan := make(chan *string)
	messageParts := []string{""}
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			messageParts = append(messageParts, text)
			if strings.Compare("\n", text) == 0 {
				output := strings.Join(messageParts, "")
				output = strings.Replace(output, "\n\n", "", 1)
				inputChan <- &output
				messageParts = []string{""}
			}
		}
	}()
	return inputChan
}

func DerpSleeper(s float32) chan int {
	timer := make(chan int)
	go func() {
		i := 0
		for {
			time.Sleep(time.Second * time.Duration(s))
			i += 1
			timer <- i
		}
	}()
	return timer
}

func main() {
	fmt.Println("Test of input")
	fmt.Println("------------------------------------------------")
	fmt.Println("Write input. Enter key on an empty line to send.")
	fmt.Println("------------------------------------------------")

	inputChan := InputGetter()
	// derpChan := DerpSleeper(10)
	for {
		prompt()
		t := <-inputChan
		fmt.Printf("Rx: %s\n", *t)
		if strings.Compare("hi", *t) == 0 {
			fmt.Println("hello, Yourself")
		}
		// select {
		// case t := <-inputChan:
		// 	fmt.Printf("Rx: %s\n", t)
		// 	prompt()
		// case <-derpChan:
		// 	fmt.Println("derp fired.")
		// }
	}
}
