package ltspice

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func parseLine(line string) string {
	split := strings.SplitN(line, ":", 2)
	if len(split) < 2 {
		log.Println("No ':' character found in string.")
		return ""
	}

	return strings.TrimSpace(split[1])
}

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

func toCSV(data map[string][]float64) (err error) {
	// buffer to write the data to
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)

	keys := make([]string, 0, len(data))
	vals := make([][]float64, 0, len(data))

	for k, v := range data {
		keys = append(keys, k)
		vals = append(vals, v)
	}

	// Write header to CSV
	if err = csvWriter.Write(keys); err != nil {
		return err
	}

	// Transpose values (columns to rows) and write to CSV
	for i := 0; i < len(vals[0]); i++ {

		row := make([]string, len(vals))
		for j := 0; j < len(vals); j++ {
			row[j] = fmt.Sprintf("%f", vals[j][i])
		}

		if err = csvWriter.Write(row); err != nil {
			return err
		}
	}

	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		return err
	}

	f, err := os.Create("myfile.csv")
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if _, err = f.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
