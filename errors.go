package ltspice

import "errors"

var (
	ErrInvalidSimulationType    = errors.New("invalid simulation type")
	ErrInvalidSimulationHeader  = errors.New("invalid simulation header")
	ErrLineTooLong              = errors.New("line too long")
	ErrInvalidUTF16             = errors.New("invalid UTF-16 sequence")
	ErrParsingError             = errors.New("parsing error")
	ErrUnexpectedEndOfFile      = errors.New("unexpected end of file")
	ErrParseStepInfo            = errors.New("parse error: failed to parse step info for stepped simulation")
	ErrTraceDoesNotExist        = errors.New("trace not found")
	ErrInvaleTraceTypeAssertion = errors.New("type assertion failed")
)
