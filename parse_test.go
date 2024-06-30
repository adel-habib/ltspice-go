package ltspice

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"math/cmplx"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadLineUTF16(t *testing.T) {
	tests := []struct {
		name    string
		input   []uint16
		want    string
		wantErr error
	}{
		{
			name:    "normal case",
			input:   []uint16{'H', 'e', 'l', 'l', 'o', '\n'},
			want:    "Hello",
			wantErr: nil,
		},
		{
			name:    "empty line",
			input:   []uint16{'\n'},
			want:    "",
			wantErr: nil,
		},
		{
			name:    "line too long",
			input:   make([]uint16, maxLineSize+1),
			want:    "",
			wantErr: ErrLineTooLong,
		},
		{
			name:    "unexpected end of file",
			input:   []uint16{'H', 'e', 'l', 'l', 'o'},
			want:    "",
			wantErr: ErrUnexpectedEndOfFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			for _, u := range tt.input {
				binary.Write(&b, binary.LittleEndian, u)
			}

			got, err := readLineUTF16(&b)

			assert.Equal(t, tt.wantErr, err, "readLineUTF16() error = %v, wantErr %v", err, tt.wantErr)
			assert.Equal(t, tt.want, got, "readLineUTF16() = %v, want %v", got, tt.want)
		})
	}
}
func TestExtractHeaderValue(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{
			name: "Normal header",
			line: "Date: Mon Mar 30 14:21:21 2020",
			want: "Mon Mar 30 14:21:21 2020",
		},
		{
			name: "Empty value",
			line: "Date:",
			want: "",
		},
		{
			name: "Leading and trailing spaces",
			line: "  Date:   Mon Mar 30 14:21:21 2020  ",
			want: "Mon Mar 30 14:21:21 2020",
		},
		{
			name: "Empty line",
			line: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractHeaderValue(tt.line)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParsingAcSimulation(t *testing.T) {

	rawFilePath := "testdata/simulations/ac/Loop-Gain/LoopGain.raw"
	resultSetPathAbs := "testdata/simulations/ac/Loop-Gain/LoopGain-abs.csv"
	resultSetPathReal := "testdata/simulations/ac/Loop-Gain/LoopGain-real.csv"
	resultSetPathImag := "testdata/simulations/ac/Loop-Gain/LoopGain-imag.csv"
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
	gotCmplx := s.complexData
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

func TestTransientSimulation(t *testing.T) {
	filePath := "testdata/simulations/trans/LM741/LM741.raw"
	resultsPath := "testdata/simulations/trans/LM741/LM741 - result-set.csv"
	s, err := Parse(filePath)
	if err != nil {
		t.Fatal(err)
	}

	resultSet, err := csvToMap(resultsPath)
	if err != nil {
		t.Fatal(err)
	}
	for k, expected := range resultSet {
		got, ok := s.data[k]
		if !ok {
			t.Fatalf("Variable %s exists in expected result set but in parsed set", k)
		}
		assert.InDeltaSlice(t, expected, got, 1e-6)
	}
}

func TestParsingNoiseSpectralSimulation(t *testing.T) {
	filePath := "testdata/simulations/noise/noise.raw"
	resultsPath := "testdata/simulations/noise/noise.csv"
	s, err := Parse(filePath)
	if err != nil {
		t.Fatal(err)
	}

	resultSet, err := csvToMap(resultsPath)
	if err != nil {
		t.Fatal(err)
	}
	for k, expected := range resultSet {
		got, ok := s.data[k]
		if !ok {
			t.Fatalf("Variable %s exists in expected result set but in parsed set", k)
		}
		assert.InDeltaSlice(t, expected, got, 1e-6)
	}
}

func TestParsingDcSweepSimulation(t *testing.T) {
	filePath := "testdata/simulations/dc/curvetrace/curvetrace.raw"
	resultsPath := "testdata/simulations/dc/curvetrace/curvetrace.csv"
	s, err := Parse(filePath)
	if err != nil {
		t.Fatal(err)
	}

	resultSet, err := csvToMap(resultsPath)
	if err != nil {
		t.Fatal(err)
	}
	for k, expected := range resultSet {
		got, ok := s.data[k]
		if !ok {
			t.Fatalf("Variable %s exists in expected result set but in parsed set", k)
		}
		assert.InDeltaSlice(t, expected, got, 1e-6)
	}
}
func TestPParsingOPSim(t *testing.T) {
	filePath := "testdata/simulations/stepped/rc/rc.raw"
	_, err := Parse(filePath)
	if err != nil {
		t.Fatal(err)
	}
}
func csvToMap(filename string) (map[string][]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	headers := lines[0]
	columns := make(map[string][]float64)

	for _, header := range headers {
		columns[header] = make([]float64, len(lines)-1)
	}

	for rowIndex, line := range lines[1:] {
		for columnIndex, cell := range line {
			value, err := strconv.ParseFloat(cell, 64)
			if err != nil {
				return nil, err
			}
			header := headers[columnIndex]
			columns[header][rowIndex] = value
		}
	}

	return columns, nil
}
