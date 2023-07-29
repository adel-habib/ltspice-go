package ltspice

import (
	"testing"
)

var result *Simulation

func BenchmarkParse(b *testing.B) {

	var sim *Simulation
	var err error

	fileName := "testdata/rc.raw"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sim, err = Parse(fileName)

		if err != nil {
			b.Fatalf("failed to parse file: %v", err)
		}
	}

	result = sim
}
