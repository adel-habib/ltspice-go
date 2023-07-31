package ltspice

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

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
	HeaderBackannotation  = "Backannotation"
)

const (
	realYAxisTraceByteSize = 4
	realXAxisTraceByteSize = 8
	complexTraceByteSize   = 16
)

type Variable struct {
	Order int // the order of the variable as it appears in the binary dataframe
	Name  string
	Typ   string // the type of the variable (time, frequency, device_voltage etc..)
	Size  int    // the size of a signle data point in bytes
}

func parseHeaderLine(r io.Reader, metadata *SimulationMetadata, line string) error {
	lineType := strings.SplitN(line, ":", 2)[0]
	switch lineType {

	case HeaderTitle:
		metadata.Title = extractHeaderValue(line)

	case HeaderPlotName:
		sim := extractHeaderValue(line)
		simType, err := SimulationTypeFromString(sim)
		if err != nil {
			return ErrInvalidSimulationType
		}
		metadata.SimType = simType

	case HeaderCommand:
		metadata.Command = extractHeaderValue(line)

	case HeaderOffset:
		num, err := strconv.ParseFloat(extractHeaderValue(line), 64)
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
		num, err := strconv.Atoi(extractHeaderValue(line))
		if err != nil {
			log.Println("Error converting string to integer:", err)
			return fmt.Errorf("%w: %s", ErrInvalidSimulationHeader, "failed to parse the number of data points from the header")
		}
		metadata.NoPoints = num

	case HeaderVariablesNumber:
		num, err := strconv.Atoi(extractHeaderValue(line))
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
			sz := realYAxisTraceByteSize
			if fields[2] == "time" || i == 0 {
				sz = realXAxisTraceByteSize
			}
			v := Variable{Order: i, Name: fields[1], Typ: fields[2], Size: sz}
			metadata.Variables[i] = v
		}
	case HeaderFlags:
		flagStr := extractHeaderValue(line)
		metadata.Flags = ParseFlags(strings.Fields(flagStr)...)
	case HeaderBackannotation:
		return nil
	default:
		log.Println("Encountered unknown header: " + lineType)

	}
	return nil
}
