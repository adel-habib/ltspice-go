package ltspice

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

const maxLineSize = 1024 // maximum number of runes in a line
const dateHeaderLayout = "Date: Mon Jan 2 15:04:05 2006"

const (
	HeaderBinary          = "binary:"
	HeaderValues          = "values:"
	HeaderTitle           = "Title"
	HeaderDate            = "Date"
	HeaderPlotName        = "Plotname"
	HeaderFlags           = "Flags"
	HeaderVariablesNumber = "No. Variables"
	HeaderPoints          = "No. Points"
	HeaderOffset          = "Offset"
	HeaderCommand         = "Command"
	HeaderVariables       = "Variables"
)

type SimulationType int

const (
	OperationPoint SimulationType = iota
	DCtransfer
	ACAanalysis
	TransientAnalysis
	NoiseSpectralDensity
	TransferFunction
)

func (s SimulationType) String() string {
	return [...]string{
		"Operation Point",
		"DC transfer characteristic",
		"AC Analysis", "Transient Analysis",
		"Noise Spectral Density",
		"Transfer Function"}[s]
}

func SimulationTypeFromString(str string) (SimulationType, error) {
	switch str {
	case "Operating Point":
		return OperationPoint, nil
	case "DC transfer characteristic":
		return DCtransfer, nil
	case "AC Analysis":
		return ACAanalysis, nil
	case "Transient Analysis":
		return TransientAnalysis, nil
	case "Noise Spectral Density":
		return NoiseSpectralDensity, nil
	case "Transfer Function":
		return TransferFunction, nil
	default:
		return 0, ErrInvalidSimulationType
	}
}

func parseHeaderLine(r io.Reader, metadata *RawFileMetadata, line string) error {
	lineType := strings.SplitN(line, ":", 2)[0]
	switch lineType {

	case HeaderTitle:
		metadata.Title = parseLine(line)

	case HeaderPlotName:
		sim := parseLine(line)
		simType, err := SimulationTypeFromString(sim)
		if err != nil {
			return ErrInvalidSimulationType
		}
		metadata.SimType = simType

	case HeaderCommand:
		metadata.Command = parseLine(line)

	case HeaderOffset:
		num, err := strconv.ParseFloat(parseLine(line), 64)
		if err != nil {
			log.Println("Error converting string to float:", err)
			metadata.Offset = 0
		} else {
			metadata.Offset = num
		}

	case HeaderDate:
		t, err := time.Parse(dateHeaderLayout, line)
		if err != nil {
			log.Println("Error parsing date:", err)
		} else {
			metadata.Date = t
		}

	case HeaderPoints:
		num, err := strconv.Atoi(parseLine(line))
		if err != nil {
			log.Println("Error converting string to integer:", err)
			return fmt.Errorf("%w: %s", ErrInvalidSimulationHeader, "failed to parse the number of data points from the header")
		}
		metadata.NoPoints = num

	case HeaderVariablesNumber:
		num, err := strconv.Atoi(parseLine(line))
		if err != nil {
			log.Println("Error converting string to integer:", err)
			return fmt.Errorf("%w: %s", ErrInvalidSimulationHeader, "failed to parse the number of variables from the header")
		}
		metadata.NoVariables = num

	case HeaderVariables:
		if metadata.NoVariables <= 0 {
			return ErrInvalidSimulationHeader
		}
		metadata.Variables = make([]Variable, metadata.NoVariables)
		for i := 0; i < metadata.NoVariables; i++ {
			l, err := readLineUTF16(r)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrInvalidSimulationHeader, err)
			}
			fields := strings.Fields(l)
			if len(fields) < 3 {
				return fmt.Errorf("%w: failed to parse variable, expected 3 fields but found %d, line: %s", ErrInvalidSimulationHeader, len(fields), l)
			}
			sz := 4
			if fields[2] == "time" || i == 0 {
				sz = 8
			}
			v := Variable{Order: i, Name: fields[1], Typ: fields[2], Size: sz}
			metadata.Variables[i] = v
		}
	case HeaderFlags:
		flagStr := parseLine(line)
		metadata.Flags = ParseFlags(strings.Fields(flagStr)...)

	default:
		log.Println("Encountered unknown header: " + lineType)

	}
	return nil
}
