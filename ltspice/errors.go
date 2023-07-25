package ltspice

import "errors"

var (
	ErrInvalidSimulationType   = errors.New("invalid simulation type")
	ErrInvalidSimulationHeader = errors.New("invalid simulation header")
	ErrLineTooLong             = errors.New("line too long")
	ErrInvalidUTF16            = errors.New("invalid UTF-16 sequence")
	ErrUnexpectedNull          = errors.New("unexpected null character")
	ErrUnexpectedError         = errors.New("unexpected error")
	ErrUnexpectedEndOfFile     = errors.New("unexpected end of file")
)
