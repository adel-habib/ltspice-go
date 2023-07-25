package ltspice

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

func parseLine(line string) string {
	split := strings.SplitN(line, ":", 2)
	if len(split) < 2 {
		fmt.Println("No ':' character found in string.")
		return ""
	}

	return strings.TrimSpace(split[1])
}

func toFloat(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
