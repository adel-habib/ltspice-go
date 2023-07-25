package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

const (
	TitleLine           = "Title"
	DateLine            = "Date"
	PlotNameLine        = "Plotname"
	FlagsLine           = "Flags"
	VariablesNumberLine = "No. Variables"
	PointsLine          = "No. Points"
	OffsetLine          = "Offset"
	CommandLine         = "Command"
	VariablesHeader     = "Variables"
	BinaryHeader        = "binary:"
)

// result can be ->
// 1- time series, a single time array with value arrays
// 2- an array of (1)
// 3- frequency analysis
// 4- OPERATION POINT

type Headers struct {
	Title             string
	Date              time.Time
	PlotName          string
	Variables         []string
	NumberOfVariables int
	NumDataPoints     int
	Offset            float64
	Command           string
	Flags             string
}

var metadata Headers

const layout = "Date: Mon Jan 2 15:04:05 2006"

func main() {
	file, err := os.Open("iter.raw")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line := readLineUTF16(reader)
		if strings.Contains(strings.ToLower(strings.TrimSpace(line)), BinaryHeader) {
			break
		}
		parseHeaderLine(reader, &metadata, line)
	}

	t := make([]float64, metadata.NumDataPoints)

out:
	for {
		timeBuffer := make([]byte, 8)
		varBuffer := make([]byte, 4)
		for dataPoint := 0; dataPoint < metadata.NumDataPoints; dataPoint++ {
			_, err := io.ReadFull(reader, timeBuffer)
			t[dataPoint] = toFloat(timeBuffer)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Reached end of file")
					break out
				}
				panic(err)
			}

			for v := 0; v < metadata.NumberOfVariables-1; v++ {
				_, err = io.ReadFull(reader, varBuffer)
				if err != nil {
					if err == io.EOF {
						fmt.Println("Reached end of file")
						break
					}
					panic(err)
				}
			}
		}
	}
	println(len(t))
}

func toFloat(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
func parseHeaderLine(r io.Reader, metadata *Headers, line string) {
	fmt.Println(line)
	lineType := strings.SplitN(line, ":", 2)[0]
	switch lineType {

	case TitleLine:
		metadata.Title = parseLine(line)

	case PlotNameLine:
		metadata.PlotName = parseLine(line)

	case CommandLine:
		metadata.Command = parseLine(line)

	case OffsetLine:
		num, err := strconv.ParseFloat(parseLine(line), 64)
		if err != nil {
			fmt.Println("Error converting string to float:", err)
			return
		}
		metadata.Offset = num

	case DateLine:
		t, err := time.Parse(layout, line)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			return
		}
		metadata.Date = t

	case PointsLine:
		num, err := strconv.Atoi(parseLine(line))
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}
		metadata.NumDataPoints = num

	case VariablesNumberLine:
		num, err := strconv.Atoi(parseLine(line))
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			return
		}
		metadata.NumberOfVariables = num

	case VariablesHeader:
		if metadata.NumberOfVariables <= 0 {
			fmt.Println(metadata.NumberOfVariables)
			panic("Invalid State")
		}
		for i := 0; i < metadata.NumberOfVariables; i++ {
			fields := strings.Fields(readLineUTF16(r))
			fmt.Println(fields)
			if len(fields) < 3 {
				panic("Failed to parse variable")
			}
			metadata.Variables = append(metadata.Variables, fields[2])
		}
	case FlagsLine:
		metadata.Flags = parseLine(line)

	default:
		fmt.Println("Invalid or unknown header: " + line)
	}
}

func parseLine(line string) string {
	split := strings.SplitN(line, ":", 2)
	if len(split) < 2 {
		fmt.Println("No ':' character found in string.")
		return ""
	}

	return strings.TrimSpace(split[1])
}

func readVariables(reader *bufio.Reader, numVars int) ([]string, error) {
	vars := make([]string, numVars)
	for i := 0; i < numVars; i++ {
		varLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading variable line: %w", err)
		}
		vars[i] = varLine
	}
	return vars, nil
}

func readLineUTF16(r io.Reader) string {
	lineBuff := make([]uint16, 0)
	for {
		runeBuff := make([]byte, 2)
		_, err := r.Read(runeBuff)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		rune := binary.LittleEndian.Uint16(runeBuff)
		if rune == '\n' {
			break
		}
		lineBuff = append(lineBuff, rune)
	}
	return string(utf16.Decode(lineBuff))
}
