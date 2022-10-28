package tapedeck

import "time"

type Message struct {
	Tick    time.Duration
	Content string
}
