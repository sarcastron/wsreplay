package tapedeck

import (
	"fmt"
	"strings"
	"time"
	"wsreplay/pkg/output"
)

// Represents the incoming messages from the websocket connection.
type Message struct {
	Tick    time.Duration
	Content []byte
}

// Interface for internal message bus. Used to decouple CLI output from business logic.
type BusMessager interface {
	// Function that outputs a message formatted for the CLI
	CliMessage() string
}

// Represents a successful or informational event.
type BusMessageInfo struct {
	Prefix  string
	Content string
}

// Generate string for CLI output
func (bm *BusMessageInfo) CliMessage() string {
	fPrefix := ""
	if bm.Prefix != "" {
		fPrefix = fmt.Sprintf("%s ", output.Info(bm.Prefix))
	}
	return fmt.Sprintf("%s%s\n", fPrefix, strings.TrimSuffix(bm.Content, "\n"))
}

// Represents an error event.
type BusMessageErr struct {
	Prefix  string
	Err     error
	IsFatal bool
}

// Generate string for CLI output
func (bm *BusMessageErr) CliMessage() string {
	fPrefix := ""
	if bm.Prefix != "" {
		fPrefix = fmt.Sprintf("%s ", output.Danger(bm.Prefix))
	}
	return fmt.Sprintf("%s%s", fPrefix, bm.Err)
}
