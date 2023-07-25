package ltspice

import (
	"log"
	"strings"
)

type Flags int

const (
	None Flags = 1 << iota
	Complex
	Forward
	Log
	Stepped
	FastAccess
)

var flagLookup = map[string]Flags{
	"real":       None,
	"complex":    Complex,
	"forward":    Forward,
	"log":        Log,
	"stepped":    Stepped,
	"fastaccess": FastAccess,
}

func ParseFlags(flagStrings ...string) Flags {
	var result Flags
	for _, flagString := range flagStrings {
		flag, ok := flagLookup[strings.ToLower(flagString)]
		if !ok {
			log.Printf("Unknown flag: %s", flagString)
			continue
		}
		result.SetFlag(flag)
	}
	return result
}

func (f Flags) HasFlag(flag Flags) bool {
	return f&flag != 0
}

func (f *Flags) SetFlag(flag Flags) {
	*f |= flag
}

func (f *Flags) ClearFlag(flag Flags) {
	*f &^= flag
}

func (f Flags) String() string {
	flagStrings := []string{}
	if f&Complex != 0 {
		flagStrings = append(flagStrings, "complex")
	}
	if f&Forward != 0 {
		flagStrings = append(flagStrings, "forward")
	}
	if f&Log != 0 {
		flagStrings = append(flagStrings, "log")
	}
	if f&Stepped != 0 {
		flagStrings = append(flagStrings, "stepped")
	}
	if f&FastAccess != 0 {
		flagStrings = append(flagStrings, "fastaccess")
	}
	return strings.Join(flagStrings, "|")
}
