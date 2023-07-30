package ltspice

import (
	"bytes"
	"encoding/binary"
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
