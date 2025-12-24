package output

import (
	"Blink/types"
	"fmt"
	"strings"
)

func ErrorOutput(err types.BlinkError) string {
	var out strings.Builder
	if err.Stage != "OK" {
		if err.Stage == "Unknown" {
			out.WriteString(fmt.Sprintf(types.Red+"[ %v ERROR ] %s\n"+types.Reset, err.Stage, err.Message))
		} else if err.Stage == "INFO" {
			out.WriteString(fmt.Sprintf(types.Yellow+"[ %v ] %v"+types.Reset, err.Stage, err.Message))
		} else {
			out.WriteString(fmt.Sprintf(types.Red+"[ %v ERROR ] %v \n"+types.Reset, err.Stage, err.Message))
		}

	}
	return out.String()
}
