package ltspice

import (
	"fmt"
	"time"
)

// SimType defines the simulation type. Supported types are Operating Point, DC transfer characteristic,
// AC Analysis, Transient Analysis, Noise Spectral Density and Transfer Function.
type SimType int

const (
	OperatingPoint SimType = iota
	DCtransfer
	ACAnalysis
	TransientAnalysis
	NoiseSpectralDensity
	TransferFunction
)

func (s SimType) String() string {
	return [...]string{
		"Operation Point",
		"DC transfer characteristic",
		"AC Analysis", "Transient Analysis",
		"Noise Spectral Density",
		"Transfer Function"}[s]
}

func simTypeFromString(str string) (SimType, error) {
	switch str {
	case "Operating Point":
		return OperatingPoint, nil
	case "DC transfer characteristic":
		return DCtransfer, nil
	case "AC Analysis":
		return ACAnalysis, nil
	case "Transient Analysis":
		return TransientAnalysis, nil
	case "Noise Spectral Density":
		return NoiseSpectralDensity, nil
	case "Noise Spectral Density - (V/Hz½ or A/Hz½)":
		return NoiseSpectralDensity, nil
	case "Transfer Function":
		return TransferFunction, nil
	default:
		return 0, ErrInvalidSimulationType
	}
}

// MetaData defines the metadata of a simulation.
type MetaData struct {
	Title        string
	Date         time.Time
	SimType      SimType
	Flags        Flags
	NoVariables  int
	NoPoints     int
	NoIterations int
	Offset       float64
	Command      string
	Variables    []Variable
	BinaryOffset int
}

// SimData defines a parsed raw file struct for an LTSpice simulation
type SimData struct {
	Meta        *MetaData
	data        map[string][]float64
	complexData map[string][]complex128
	xAxisLabel  string
	steps       *steps
	stepPoints  int
}

// GetType retrieves the type of the simulation
func (sim *SimData) GetType() SimType {
	return sim.Meta.SimType
}

// GetSteps retrieves the number of steps in the simulation
func (sim *SimData) GetSteps() int {
	return sim.steps.count
}

func (sim *SimData) GetVariables() []Variable {
	return sim.Meta.Variables
}

// Trace represents the simulation data for a variable in an LTSpice simulation.
// It contains the variable name and data points, which can be either float64 for regular simulations
// or complex128 for AC analysis. For stepped simulations, the data includes a signal per each step.
// The Data field is a flat slice that contains all signals.
type Trace[T float64 | complex128] struct {
	s    *steps
	Name string
	Data []T
}

// GetTrace retrieves the trace with the specified name from the given SimData.
// The type of the trace T must be specified, which can be either float64 or complex128, depending on the simulation type.
//
// The SimData.Meta.Flags field indicates the nature of the data, such as whether it is complex (Complex flag).
// The SimData.GetType function returns the type of simulation (e.g., AC Analysis, Transient Analysis).
//
// If the traceName does not exist, ErrTraceDoesNotExist is returned.
// If the type assertion fails, ErrInvalidTraceTypeAssertion is returned.
//
// Example usage:
//
//	trace, err := ltspice.GetTrace[complex128](simData, "V(out)")
//	if err != nil {
//	    log.Fatalf("Failed to get trace: %v", err)
//	}
//
//	signal := trace.GetSignal()
//	fmt.Printf("Signal data: %v\n", signal)
func GetTrace[T float64 | complex128](sim *SimData, name string) (*Trace[T], error) {
	var traceData []T

	if sim.Meta.Flags.hasFlag(Complex) {
		if data, ok := sim.complexData[name]; ok {
			if complexData, ok := any(data).([]T); ok {
				traceData = complexData
			} else {
				return nil, fmt.Errorf("type assertion to complex128 on trace %s failed %w", name, ErrInvaleTraceTypeAssertion)
			}
		} else {
			return nil, ErrTraceDoesNotExist
		}
	} else {
		if data, ok := sim.data[name]; ok {
			if floatData, ok := any(data).([]T); ok {
				traceData = floatData
			} else {
				return nil, fmt.Errorf("type assertion to float64 on trace %s failed %w", name, ErrInvaleTraceTypeAssertion)
			}
		} else {
			return nil, ErrTraceDoesNotExist
		}
	}

	return &Trace[T]{Name: name, Data: traceData, s: sim.steps}, nil
}

// GetSignal returns the data contained in the trace.
//
// For stepped simulations, a step index should be passed to retrieve data for a specific step.
// The number of steps in the simulation can be determined using the GetSteps() method on SimData.
//
// - If no step index is provided, the data for the first step is returned.
// - If an invalid step index is provided (out of range), an empty slice is returned.
// - If the simulation is not stepped, all the data is returned regardless of any step index passed.
//
// Example usage:
//
//	// Retrieve the trace for the current signal
//	IC2, err := ltspice.GetTrace[float64](simData, "I(C2)")
//	if err != nil {
//	    log.Fatalf("Failed to get trace: %v", err)
//	}
//
//	// Determine the number of steps in the simulation
//	steps := simData.GetSteps()
//
//	// Plot the signal for each step
//	for step := 0; step < steps; step++ {
//	    timeWave := simData.GetXAxis(step)
//	    currentWave := IC2.GetSignal(step)
//	    // Assuming a plotting function
//	    plot(timeWave, currentWave, fmt.Sprintf("Step %d", step))
//	}
func (t *Trace[T]) GetSignal(step ...int) []T {
	if t.s.count <= 1 {
		// If the simulation is not stepped, return all data
		return t.Data
	}

	if len(step) == 0 {
		// If no step index is provided, return the data for the first step
		return t.Data[0:t.s.offsets[1]]
	}

	stepIndex := step[0]
	if stepIndex < 0 || stepIndex >= t.s.count {
		// If the step index is out of range, return an empty slice
		return []T{}
	}

	start := t.s.offsets[stepIndex]
	end := len(t.Data)
	if stepIndex+1 < len(t.s.offsets) {
		end = t.s.offsets[stepIndex+1]
	}

	return t.Data[start:end]
}

// GetXAxis returns the x-axis data for the simulation.
// For stepped simulations, a step index can be passed to retrieve the x-axis data for a specific step.
// If no step index is provided, the data for the first step is returned.
// If anything goes wrong, an empty slice is returned.
//
// Example usage:
//
//	xAxis := simData.GetXAxis()
//	xAxisStep1 := simData.GetXAxis(1)
func (sim *SimData) GetXAxis(step ...int) []float64 {
	// Determine the step index
	stepIndex := 0
	if len(step) > 0 {
		stepIndex = step[0]
	}

	// Check if the simulation is complex (e.g., AC analysis)
	if sim.Meta.Flags.hasFlag(Complex) {
		xTrace, err := GetTrace[complex128](sim, sim.xAxisLabel)
		if err != nil {
			return []float64{}
		}

		complexData := xTrace.GetSignal(stepIndex)
		xAxis := make([]float64, len(complexData))
		for i, c := range complexData {
			xAxis[i] = real(c)
		}
		return xAxis
	}

	// Retrieve the trace for the time axis (or equivalent)
	xTrace, err := GetTrace[float64](sim, sim.xAxisLabel)
	if err != nil {
		return []float64{}
	}

	xAxis := xTrace.GetSignal(stepIndex)
	return xAxis
}
