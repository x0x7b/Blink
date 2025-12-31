package output

import (
	"Blink/types"
	"fmt"
	"strings"
	"time"
)

var lastRender time.Time
var start = time.Now()

func Report(p types.Progress) {
	RenderProgress(p)
}

func RenderProgress(p types.Progress) {
	if p.Current == 1 {
		start = time.Now()
	}

	if time.Since(lastRender) < 200*time.Millisecond {
		return
	}
	lastRender = time.Now()

	percent := float64(p.Current) / float64(p.Total) * 100
	elapsed := time.Since(start).Seconds()

	if elapsed <= 0 {
		elapsed = 0.001
	}

	speed := float64(p.Current) / elapsed

	barWidth := 24
	filled := int(percent / 100 * float64(barWidth))
	fmt.Printf(
		"\r\033[K[%s:%s] [%s%s] %5.2f%% | %.1f r/s",
		p.Stage,
		p.Target,
		strings.Repeat("█", filled),
		strings.Repeat("░", barWidth-filled),
		percent,
		speed,
	)
	if percent == 100 {
		fmt.Print("\r\033[K\n")
	}
}
