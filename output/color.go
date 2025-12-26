package output

import (
	"Blink/types"
	"time"
)

func ColorStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return types.Green
	case code >= 300 && code < 400:
		return types.Blue
	case code >= 400 && code < 500:
		return types.Yellow
	default:
		return types.Red
	}
}

func colorTime(timing time.Duration) string {
	switch {
	case timing < 150*time.Millisecond:
		return types.Green
	case timing >= 150*time.Millisecond && timing < 700*time.Millisecond:
		return types.Yellow
	case timing >= 700*time.Millisecond:
		return types.Red
	default:
		return types.Red
	}

}

func colorScore(score float64) string {
	switch {
	case score <= 0.3:
		return types.White
	case score > 0.3 && score < 0.7:
		return types.Yellow
	case score >= 0.7:
		return types.Red
	default:
		return types.White
	}
}
