package tapedeck

import "time"

// TODO switch this to bytes. No need to convert it.
type Message struct {
	Tick    time.Duration
	Content []byte
}
