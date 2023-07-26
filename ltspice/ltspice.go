package ltspice

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"unicode/utf16"
)

func Parse(fileName string) (*RawFileMetadata, error) {
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
	parseBinaryData(reader, meta)
	return meta, nil
}

func readLineUTF16(r io.Reader) (string, error) {
	lineBuff := make([]uint16, 0)
	for {
		if len(lineBuff) > maxLineSize {
			return "", ErrLineTooLong
		}

		var rune uint16
		err := binary.Read(r, binary.LittleEndian, &rune)

		if err != nil {
			if errors.Is(err, io.EOF) {
				return "", ErrUnexpectedEndOfFile
			} else {
				return "", ErrUnexpectedError
			}
		}

		if rune == '\n' {
			fmt.Println(string(utf16.Decode(lineBuff)))
			return string(utf16.Decode(lineBuff)), nil
		}

		lineBuff = append(lineBuff, rune)
	}
}

func parseHeaders(reader io.Reader) (*RawFileMetadata, error) {
	var metadata = &RawFileMetadata{Flags: None}
	for {
		line, err := readLineUTF16(reader)
		if err != nil {
			return nil, err
		}
		if strings.Contains(strings.ToLower(strings.TrimSpace(line)), HeaderBinary) || strings.Contains(strings.ToLower(strings.TrimSpace(line)), HeaderValues) {
			break
		}
		err = parseHeaderLine(reader, metadata, line)
		if err != nil {
			return nil, err
		}
	}
	return metadata, nil
}

func parseBinaryData(reader io.Reader, meta *RawFileMetadata) error {
	// data := make([][]float64, meta.NoVariables)
	for i := 0; i < meta.NoPoints; i++ {
		for j := 0; j < meta.NoVariables; j++ {
			v := meta.Variables[j]
			buffer := make([]byte, v.Size)
			n, err := io.ReadFull(reader, buffer)
			if err != nil {
				return err

			}
			if v.Size == 4 {
				bits := binary.LittleEndian.Uint32(buffer)
				float := math.Float32frombits(bits)
				fmt.Printf("%f\n", float)
			} else {
				float := toFloat(buffer)
				fmt.Printf("%f\n", float)
			}

		}
	}

	return nil
}
