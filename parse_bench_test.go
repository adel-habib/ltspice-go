package main

import (
	"testing"

	"github.com/theadell/ltspice-go/ltspice"
)

var result *ltspice.Simulation

func BenchmarkParse(b *testing.B) {

	var sim *ltspice.Simulation
	var err error

	fileName := "rc.raw"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sim, err = ltspice.Parse(fileName)

		if err != nil {
			b.Fatalf("failed to parse file: %v", err)
		}
	}

	result = sim
}
