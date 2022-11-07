package output

import "time"

type Spinner struct {
	Frames       []string
	tickDuration time.Duration
	Index        int
}

// Starts the spinner.
func (s *Spinner) Spin() chan string {
	spinnerChan := make(chan string)
	go func() {
		ticker := time.NewTicker(s.tickDuration)
		defer ticker.Stop()
		for {
			<-ticker.C
			s.Tick()
		}
	}()
	return spinnerChan
}

// Return string representation of the current state of the spinner.
func (s *Spinner) Render() string {
	return string(s.Frames[s.Index])
}

// Advances the spinner one frame. Used by Spin method.
func (s *Spinner) Tick() {
	if s.Index >= len(s.Frames)-1 {
		s.Index = 0
	} else {
		s.Index += 1
	}
}

var PlaybackSpinner = &Spinner{
	Frames:       []string{"⣾", "⣷", "⣯", "⣟", "⡿", "⢿", "⣻", "⣽"},
	tickDuration: time.Millisecond * 200,
}

var MeterSpinner = &Spinner{
	Frames: []string{
		"▱▱▱",
		"▰▱▱",
		"▰▰▱",
		"▰▰▰",
		"▰▰▱",
		"▰▱▱",
		"▱▱▱",
	},
	tickDuration: time.Millisecond * 150,
}
