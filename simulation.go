package ltspice

import "time"

type SimulationType int

const (
	OperatingPoint SimulationType = iota
	DCtransfer
	ACAnalysis
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

type SimulationMetadata struct {
	Title        string
	Date         time.Time
	SimType      SimulationType
	Flags        Flags
	NoVariables  int
	NoPoints     int
	Offset       float64
	Command      string
	Traces       []Trace
	BinaryOffset int
}

type Simulation struct {
	metaData    *SimulationMetadata
	Data        map[string][]float64
	ComplexData map[string][]complex128
}

func (sim *Simulation) GetType() SimulationType {
	return sim.metaData.SimType
}

func (sim *Simulation) GetNumberOfDataPoints() int {
	return sim.metaData.NoPoints
}

func (sim *Simulation) GetNumberOfVariables() int {
	return sim.metaData.NoVariables
}

func (sim *Simulation) GetVariables() []Trace {
	return sim.metaData.Traces
}
func (sim *Simulation) GetVariableNames() []string {
	names := make([]string, sim.metaData.NoVariables)
	for i, v := range sim.metaData.Traces {
		names[i] = v.Name
	}
	return names
}

func (sim *Simulation) GetFlags() Flags {
	return sim.metaData.Flags
}
