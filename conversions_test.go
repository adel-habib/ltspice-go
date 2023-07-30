package ltspice

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToFloat(t *testing.T) {
	floatsToTest := []float64{
		0,
		math.SmallestNonzeroFloat64,
		0.000000000000001,
		0.123456789,
		1,
		10,
		1000.5,
		math.MaxFloat64,
		-math.MaxFloat64,
	}

	for _, originalFloat := range floatsToTest {
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, originalFloat)
		assert.NoError(t, err, "Failed to write float64 to buffer")

		recoveredFloat := toFloat(buf.Bytes())
		assert.Equal(t, originalFloat, recoveredFloat, "toFloat did not correctly recover the original float64 value")
	}
}

func TestToFloatFrom32(t *testing.T) {
	// Prepare a list of floats to test
	floatsToTest := []float32{
		0,
		math.SmallestNonzeroFloat32,
		0.0000000001,
		0.123456,
		1,
		10,
		1000.5,
		math.MaxFloat32,
		-math.MaxFloat32,
	}

	for _, originalFloat := range floatsToTest {
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, originalFloat)
		assert.NoError(t, err, "Failed to write float32 to buffer")

		recoveredFloat := toFloatFrom32(buf.Bytes())
		assert.Equal(t, float64(originalFloat), recoveredFloat, "toFloatFrom32 did not correctly recover the original float32 value")
	}
}
