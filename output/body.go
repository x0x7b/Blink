package output

import "Blink/types"

func bodyOutput(bl types.BlinkResponse, fc types.FlagCondition) string {
	if fc.ShowBody {
		return bl.BodyPreview
	} else if fc.ShowFullBody {
		return string(bl.Body)
	}
	return ""
}
