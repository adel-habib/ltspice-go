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
	headerBinary          = "binary:"
	headerValues          = "values:"
	headerTitle           = "Title"
	headerDate            = "Date"
	headerPlotName        = "Plotname"
	headerFlags           = "Flags"
	headerVariablesNumber = "No. Variables"
	headerPoints          = "No. Points"
	headerOffset          = "Offset"
	headerCommand         = "Command"
	headerVariables       = "Variables"
	headerBackannotation  = "Backannotation"
)

const (
	realYAxisTraceByteSize = 4
	realXAxisTraceByteSize = 8
	complexTraceByteSize   = 16
)

type Variable struct {
	order int // the order of the variable as it appears in the binary dataframe
	Name  string
	Typ   string // the type of the variable (time, frequency, device_voltage etc..)
	size  int    // the size of a signle data point in bytes
}

func parseHeaderLine(r io.Reader, metadata *MetaData, line string) error {
	lineType := strings.SplitN(line, ":", 2)[0]
	switch lineType {

	case headerTitle:
		metadata.Title = extractHeaderValue(line)

	case headerPlotName:
		sim := extractHeaderValue(line)
		simType, err := simTypeFromString(sim)
		if err != nil {
			return ErrInvalidSimulationType
		}
		metadata.SimType = simType

	case headerCommand:
		metadata.Command = extractHeaderValue(line)

	case headerOffset:
		num, err := strconv.ParseFloat(extractHeaderValue(line), 64)
		if err != nil {
			log.Println("Error converting string to float:", err)
			metadata.Offset = 0
		} else {
			metadata.Offset = num
		}

	case headerDate:
		t, err := time.Parse(dateHeaderLayout, line)
		if err != nil {
			log.Println("Error parsing date:", err)
		} else {
			metadata.Date = t
		}

	case headerPoints:
		num, err := strconv.Atoi(extractHeaderValue(line))
		if err != nil {
			log.Println("Error converting string to integer:", err)
			return fmt.Errorf("%w: %s", ErrInvalidSimulationHeader, "failed to parse the number of data points from the header")
		}
		metadata.NoPoints = num

	case headerVariablesNumber:
		num, err := strconv.Atoi(extractHeaderValue(line))
		if err != nil {
			log.Println("Error converting string to integer:", err)
			return fmt.Errorf("%w: %s", ErrInvalidSimulationHeader, "failed to parse the number of variables from the header")
		}
		metadata.NoVariables = num

	case headerVariables:
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
			v := Variable{order: i, Name: fields[1], Typ: fields[2], size: sz}
			metadata.Variables[i] = v
		}
	case headerFlags:
		flagStr := extractHeaderValue(line)
		metadata.Flags = parseFlags(strings.Fields(flagStr)...)
	case headerBackannotation:
		return nil
	default:
		log.Println("Encountered unknown header: " + lineType)

	}
	return nil
}
