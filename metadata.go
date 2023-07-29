package ltspice

import (
	"fmt"
	"strings"
	"time"
)

type RawFileMetadata struct {
	Title        string
	Date         time.Time
	SimType      SimulationType
	Flags        Flags
	NoVariables  int
	NoPoints     int
	Offset       float64
	Command      string
	Variables    []Variable
	BinaryOffset int
}

func (rfm RawFileMetadata) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Title: %s\n", rfm.Title))
	builder.WriteString(fmt.Sprintf("Date: %s\n", rfm.Date.Format(time.ANSIC)))
	builder.WriteString(fmt.Sprintf("Plotname: %s\n", rfm.SimType.String()))
	builder.WriteString(fmt.Sprintf("Flags: %s\n", rfm.Flags.String()))
	builder.WriteString(fmt.Sprintf("No. Variables: %d\n", rfm.NoVariables))
	builder.WriteString(fmt.Sprintf("No. Points: %d\n", rfm.NoPoints))
	builder.WriteString(fmt.Sprintf("Offset: %.2f\n", rfm.Offset))
	builder.WriteString(fmt.Sprintf("Command: %s\n", rfm.Command))

	builder.WriteString("Variables:\n")
	for _, v := range rfm.Variables {
		builder.WriteString(fmt.Sprintf("\t%d\t%v\n", v.Order, v))
	}

	return builder.String()
}
