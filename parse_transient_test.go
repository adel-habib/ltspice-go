package ltspice

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransientSimulation(t *testing.T) {
	filePath := "testdata/simulations/LM741/LM741.raw"
	resultsPath := "testdata/simulations/LM741/LM741 - result-set.csv"
	s, err := Parse(filePath)
	if err != nil {
		t.Fatal(err)
	}

	resultSet, err := csvToMap(resultsPath)
	if err != nil {
		t.Fatal(err)
	}
	for k, expected := range resultSet {
		got, ok := s.Data[k]
		if !ok {
			t.Fatalf("Variable %s exists in expected result set but in parsed set", k)
		}
		assert.InDeltaSlice(t, expected, got, 1e-4)
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
