package ltspice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlags(t *testing.T) {
	t.Run("ParseFlags", func(t *testing.T) {
		tests := []struct {
			name  string
			input []string
			want  Flags
		}{
			{
				name:  "No Flags",
				input: []string{},
				want:  0,
			},
			{
				name:  "Single Flag",
				input: []string{"complex"},
				want:  Complex,
			},
			{
				name:  "Multiple Flags",
				input: []string{"complex", "forward", "log"},
				want:  Complex | Forward | Log,
			},
			{
				name:  "Unknown Flag",
				input: []string{"complex", "unknown"},
				want:  Complex,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got := ParseFlags(tc.input...)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("Set and Clear Flags", func(t *testing.T) {
		f := None
		f.SetFlag(Complex)
		assert.True(t, f.HasFlag(Complex))
		f.ClearFlag(Complex)
		assert.False(t, f.HasFlag(Complex))
	})

	t.Run("String", func(t *testing.T) {
		tests := []struct {
			name string
			f    Flags
			want string
		}{
			{
				name: "No Flags",
				f:    None,
				want: "",
			},
			{
				name: "Single Flag",
				f:    Complex,
				want: "complex",
			},
			{
				name: "Multiple Flags",
				f:    Complex | Forward | Log,
				want: "complex|forward|log",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				got := tc.f.String()
				assert.Equal(t, tc.want, got)
			})
		}
	})
}
