package output

import (
	"Blink/types"
	"fmt"
	"strings"
	"time"
)

func Report(p types.Progress) {
	RenderProgress(p)
}

func RenderProgress(p types.Progress) {
	var lastRender time.Time
	var start = time.Now()
	if time.Since(lastRender) < 200*time.Millisecond {
		return
	}
	lastRender = time.Now()

	percent := float64(p.Current) / float64(p.Total) * 100
	elapsed := time.Since(start).Seconds()
	speed := float64(p.Current) / elapsed

	barWidth := 24
	filled := int(percent / 100 * float64(barWidth))
	fmt.Printf(
		"\r\033[K[%s:%s] [%s%s] %5.2f%% | %.1fk/s",
		p.Stage,
		p.Target,
		strings.Repeat("█", filled),
		strings.Repeat("░", barWidth-filled),
		percent,
		speed/1000,
	)
	if percent == 100 {
		fmt.Print("\r\033[K")
	}
}
