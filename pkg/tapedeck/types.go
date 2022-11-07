package tapedeck

import (
	"fmt"
	"strings"
	"time"
	"wsreplay/pkg/output"
)

// TODO switch this to bytes. No need to convert it.
type Message struct {
	Tick    time.Duration
	Content []byte
}

type BusMessager interface {
	// Function that outputs a message formatted for the CLI
	CliMessage() string
}

type BusMessageInfo struct {
	Prefix  string
	Content string
}

func (bm *BusMessageInfo) CliMessage() string {
	fPrefix := ""
	if bm.Prefix != "" {
		fPrefix = fmt.Sprintf("%s ", output.Info(bm.Prefix))
	}
	return fmt.Sprintf("%s%s\n", fPrefix, strings.TrimSuffix(bm.Content, "\n"))
}

type BusMessageErr struct {
	Prefix  string
	Err     error
	IsFatal bool
}

func (bm *BusMessageErr) CliMessage() string {
	fPrefix := ""
	if bm.Prefix != "" {
		fPrefix = fmt.Sprintf("%s ", output.Danger(bm.Prefix))
	}
	return fmt.Sprintf("%s%s", fPrefix, bm.Err)
}
