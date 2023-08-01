package ltspice

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf16"
)

const maxLineSize = 1024
const dateHeaderLayout = "Date: Mon Jan 2 15:04:05 2006"

func Parse(fileName string) (*Simulation, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	meta, err := parseHeaders(reader)
	if err != nil {
		return nil, err
	}
	sim := &Simulation{
		metaData: meta,
	}
	if !meta.Flags.HasFlag(Complex) {
		data, err := parseBinaryData(reader, meta)
		if err != nil {
			return nil, err
		}
		sim.Data = data

		return sim, nil
	} else {
		data, err := parseBinaryComplex(reader, meta)
		if err != nil {
			return nil, err
		}
		sim.ComplexData = data
		return sim, nil
	}

}

func parseHeaders(reader io.Reader) (*SimulationMetadata, error) {
	var metadata = &SimulationMetadata{Flags: None}
	for {
		line, err := readLineUTF16(reader)
		if err != nil {
			return nil, err
		}
		if strings.Contains(strings.ToLower(strings.TrimSpace(line)), HeaderBinary) || strings.Contains(strings.ToLower(strings.TrimSpace(line)), HeaderValues) {
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

func parseBinaryData(reader io.Reader, meta *SimulationMetadata) (map[string][]float64, error) {
	data := make(map[string][]float64)
	for _, v := range meta.Traces {
		data[v.Name] = make([]float64, meta.NoPoints)
	}
	buff := make([]byte, 16)
	for i := 0; i < meta.NoPoints; i++ {
		for _, v := range meta.Traces {
			_, err := io.ReadFull(reader, buff[:v.Size])
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			var val float64
			if v.Size == 4 {
				val = toFloatFrom32(buff[:v.Size])
			} else {
				val = toFloat(buff[:v.Size])
			}
			data[v.Name][i] = val
		}
	}
	return data, nil
}

func parseBinaryComplex(reader io.Reader, meta *SimulationMetadata) (map[string][]complex128, error) {
	data := make(map[string][]complex128)
	for _, v := range meta.Traces {
		data[v.Name] = make([]complex128, meta.NoPoints)
	}
	buff := make([]byte, 16)
	for i := 0; i < meta.NoPoints; i++ {
		for _, v := range meta.Traces {
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
				return "", ErrUnexpectedError
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
