package ltspice

import (
	"encoding/binary"
	"math"
)

func toFloat(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
func toFloatFrom32(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint32(bytes)
	n := math.Float32frombits(bits)
	return float64(n)
}
