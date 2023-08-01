package ltspice

import (
	"math/cmplx"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingAcSimulation(t *testing.T) {

	rawFilePath := "testdata/simulations/AC-Loop-Gain/LoopGain.raw"
	resultSetPathAbs := "testdata/simulations/AC-Loop-Gain/LoopGain-abs.csv"
	resultSetPathReal := "testdata/simulations/AC-Loop-Gain/LoopGain-real.csv"
	resultSetPathImag := "testdata/simulations/AC-Loop-Gain/LoopGain-imag.csv"
	s, err := Parse(rawFilePath)
	if err != nil {
		t.Fail()
	}
	resultSetAbs, err := csvToMap(resultSetPathAbs)
	if err != nil {
		t.Fatal(err)
	}
	resultSetReal, err := csvToMap(resultSetPathReal)
	if err != nil {
		t.Fatal(err)
	}
	resultSetImag, err := csvToMap(resultSetPathImag)
	if err != nil {
		t.Fatal(err)
	}
	gotCmplx := s.ComplexData
	gotAbs := make(map[string][]float64)
	gotReal := make(map[string][]float64)
	gotImag := make(map[string][]float64)
	for k, v := range gotCmplx {
		gotAbs[k] = make([]float64, len(v))
		gotReal[k] = make([]float64, len(v))
		gotImag[k] = make([]float64, len(v))
		for idx, item := range v {
			gotAbs[k][idx] = cmplx.Abs(item)
			gotReal[k][idx] = real(item)
			gotImag[k][idx] = imag(item)
		}
	}

	for k := range gotCmplx {
		expectedAbs, ok := resultSetAbs[k]
		if !ok {
			t.Fatalf("Variable %s exists in expected result set but in parsed set", k)
		}
		assert.InDeltaSlice(t, expectedAbs, gotAbs[k], 1e-6)

		expectedReal, ok := resultSetReal[k]
		if !ok {
			t.Fatalf("Variable %s exists in expected result set but in parsed set", k)
		}
		assert.InDeltaSlice(t, expectedReal, gotReal[k], 1e-6)

		expectedImag, ok := resultSetImag[k]
		if !ok {
			t.Fatalf("Variable %s exists in expected result set but in parsed set", k)
		}
		assert.InDeltaSlice(t, expectedImag, gotImag[k], 1e-6)
	}

}
