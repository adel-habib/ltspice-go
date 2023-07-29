package ltspice

type Variable struct {
	Order int // the order of the variable as it appears in the binary dataframe
	Name  string
	Typ   string // the type of the variable (time, frequency, device_voltage etc..)
	Size  int    // the size of a signle data point in bytes
}
