package ltspice

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strings"
	"unicode/utf16"
)

func Parse(fileName string) (*RawFileMetadata, error) {
	var metadata = &RawFileMetadata{Flags: None}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

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

func readLineUTF16(r io.Reader) (string, error) {
	lineBuff := make([]uint16, 0)
	for {
		if len(lineBuff) > maxLineSize {
			return "", ErrLineTooLong
		}

		// Read 2 bytes and interpret them as UTF-16 Rune
		// runeBuff := make([]byte, 2)
		// l, err := r.Read(runeBuff)
		// if err != nil {
		//	return "", err
		// }
		// if l < 2 {
		// 	continue
		// }
		// rune := binary.LittleEndian.Uint16(runeBuff)
		// if rune == '\n' {
		//	break
		// }

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
			return string(utf16.Decode(lineBuff)), nil
		}

		lineBuff = append(lineBuff, rune)
	}
}
