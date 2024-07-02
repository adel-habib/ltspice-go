package ltspice

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"strings"
	"unicode/utf16"
)

const maxLineSize = 1024
const dateHeaderLayout = "Date: Mon Jan 2 15:04:05 2006"

// Parse loads and parses an LTSpice raw data file specified by fileName.
// It returns the parsed simulation data as a SimData object which can be used to access the data inside the
// RAW file.
// If an error occurs during parsing, it returns a non-nil error.
//
// Example usage:
//
//	simData, err := ltspice.Parse("path/to/ltspice.raw")
//	if err != nil {
//	    log.Fatalf("Failed to parse LTSpice raw data: %v", err)
//	}
func Parse(fileName string) (*SimData, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	return ParseFromReader(reader)
}

// ParseFromReader parses LTSpice raw data file from the provided io.Reader.
// It returns the parsed simulation data as a SimData object, or a non-nil error if an error occurs.
func ParseFromReader(reader io.Reader) (*SimData, error) {
	meta, err := parseHeaders(reader)
	if err != nil {
		return nil, err
	}

	toReadBytes := 0
	for _, t := range meta.Variables {
		toReadBytes += (t.size * meta.NoPoints)
	}

	binaryData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(binaryData)

	sim := &SimData{
		Meta:       meta,
		xAxisLabel: meta.Variables[0].Name,
	}
	if !meta.Flags.hasFlag(Complex) {
		data, err := parseBinaryData(r, meta)
		if err != nil {
			return nil, err
		}
		sim.data = data
	} else {
		data, err := parseBinaryComplex(r, meta)
		if err != nil {
			return nil, err
		}
		sim.complexData = data
	}
	steps := &steps{
		count:   meta.NoPoints,
		offsets: make([]int, 0),
	}
	sim.stepPoints = sim.Meta.NoPoints
	if sim.Meta.Flags.hasFlag(Stepped) {
		switch sim.Meta.SimType {
		case OperatingPoint, TransferFunction:
			steps.count = len(sim.data[sim.GetVariables()[0].Name])
		default:
			var xAxis []float64
			if sim.Meta.Flags.hasFlag(Complex) {
				for _, c := range sim.complexData[sim.xAxisLabel] {
					xAxis = append(xAxis, real(c))
				}
			} else {
				xAxis = sim.data[sim.xAxisLabel]
			}
			steps, err = detectSteps(xAxis)
			if err != nil {
				return nil, err
			}
		}
		sim.steps = steps
	}
	return sim, nil

}

func parseHeaders(reader io.Reader) (*MetaData, error) {
	var metadata = &MetaData{Flags: None}
	for {
		line, err := readLineUTF16(reader)
		if err != nil {
			return nil, err
		}
		if strings.Contains(strings.ToLower(strings.TrimSpace(line)), headerBinary) || strings.Contains(strings.ToLower(strings.TrimSpace(line)), headerValues) {
			break
		}
		if line == "" {
			continue
		}
		err = parseHeaderLine(reader, metadata, line)
		if err != nil {
			return nil, err
		}
	}
	return metadata, nil
}

func parseBinaryData(reader io.Reader, meta *MetaData) (map[string][]float64, error) {
	data := make(map[string][]float64)
	for _, v := range meta.Variables {
		data[v.Name] = make([]float64, meta.NoPoints)
	}
	buff := make([]byte, 16)
	for i := 0; i < meta.NoPoints; i++ {
		for _, v := range meta.Variables {
			_, err := io.ReadFull(reader, buff[:v.size])
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			var val float64
			if v.size == 4 {
				val = toFloatFrom32(buff[:v.size])
			} else {
				val = toFloat(buff[:v.size])
			}
			data[v.Name][i] = val
		}
	}
	return data, nil
}

func parseBinaryComplex(reader io.Reader, meta *MetaData) (map[string][]complex128, error) {
	data := make(map[string][]complex128)
	for _, v := range meta.Variables {
		data[v.Name] = make([]complex128, meta.NoPoints)
	}
	buff := make([]byte, 16)
	for i := 0; i < meta.NoPoints; i++ {
		for _, v := range meta.Variables {
			_, err := io.ReadFull(reader, buff[:16])
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			var real float64
			var img float64
			real = toFloat(buff[:8])
			img = toFloat(buff[8:16])
			data[v.Name][i] = complex(real, img)
		}
	}
	return data, nil
}
func readLineUTF16(r io.Reader) (string, error) {
	lineBuff := make([]uint16, 0, maxLineSize)
	buff := make([]byte, 2)
	for {
		if len(lineBuff) > maxLineSize {
			return "", ErrLineTooLong
		}

		_, err := io.ReadFull(r, buff)

		if err != nil {
			if errors.Is(err, io.EOF) {
				return "", ErrUnexpectedEndOfFile
			} else {
				return "", ErrParsingError
			}
		}
		rune := binary.LittleEndian.Uint16(buff)
		if rune == '\n' {
			return strings.TrimSpace(string(utf16.Decode(lineBuff))), nil
		}

		lineBuff = append(lineBuff, rune)
	}
}

func extractHeaderValue(line string) string {
	split := strings.SplitN(line, ":", 2)
	if len(split) < 2 {
		return ""
	}

	return strings.TrimSpace(split[1])
}

type steps struct {
	count   int
	offsets []int
}

func detectSteps(xAxis []float64) (*steps, error) {
	var s = &steps{
		count:   0,
		offsets: make([]int, 0),
	}
	if len(xAxis) == 0 {
		return nil, ErrParseStepInfo
	}

	origin := xAxis[0]
	for idx, point := range xAxis {
		if math.Abs(float64(point-origin)) < 1e-10 {
			s.count += 1
			s.offsets = append(s.offsets, idx)
		}
	}

	if s.count == 0 {
		slog.Error("failed to detect steps or find a pattern in x-axis data")
		return nil, ErrParseStepInfo
	}
	return s, nil
}
