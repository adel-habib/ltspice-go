package ltspice

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"
	"unicode/utf16"

	"github.com/stretchr/testify/assert"
)

func TestSimulationType_String(t *testing.T) {
	tests := []struct {
		name string
		s    SimType
		want string
	}{
		{
			name: "Operation Point",
			s:    OperatingPoint,
			want: "Operation Point",
		},
		{
			name: "DC transfer characteristic",
			s:    DCtransfer,
			want: "DC transfer characteristic",
		},
		{
			name: "AC Analysis",
			s:    ACAnalysis,
			want: "AC Analysis",
		},
		{
			name: "Transient Analysis",
			s:    TransientAnalysis,
			want: "Transient Analysis",
		},
		{
			name: "Noise Spectral Density",
			s:    NoiseSpectralDensity,
			want: "Noise Spectral Density",
		},
		{
			name: "Transfer Function",
			s:    TransferFunction,
			want: "Transfer Function",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSimulationTypeFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    SimType
		wantErr error
	}{
		{
			name:  "Operation Point",
			input: "Operating Point",
			want:  OperatingPoint,
		},
		{
			name:  "DC transfer characteristic",
			input: "DC transfer characteristic",
			want:  DCtransfer,
		},
		{
			name:  "AC Analysis",
			input: "AC Analysis",
			want:  ACAnalysis,
		},
		{
			name:  "Transient Analysis",
			input: "Transient Analysis",
			want:  TransientAnalysis,
		},
		{
			name:  "Noise Spectral Density",
			input: "Noise Spectral Density",
			want:  NoiseSpectralDensity,
		},
		{
			name:  "Transfer Function",
			input: "Transfer Function",
			want:  TransferFunction,
		},
		{
			name:    "Invalid",
			input:   "Invalid",
			want:    0,
			wantErr: ErrInvalidSimulationType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := simTypeFromString(tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestParseHeaders(t *testing.T) {
	test_1_headers := `Title: * Z:\home\wine\ltspice\Draft1.asc
						Date: Tue Jul 25 12:15:28 2023
						Plotname: Transient Analysis
						Flags: real forward stepped
						No. Variables: 4
						No. Points:          142
						Offset:   0.0000000000000000e+000
						Command: Linear Technology Corporation LTspice XVII
						Variables:
										0   time    time
										1   V(n001) voltage
										2   I(R1)   device_current
										3   I(V1)   device_current
						Binary: 
`

	test_2_headers := `Title: * Z:\home\wine\ltspice\test.asc
						Date: Sat Jul 29 12:38:42 2023
						Plotname: Transient Analysis
						Flags: real forward
						No. Variables: 17
						No. Points:          520
						Offset:   0.0000000000000000e+000
						Command: Linear Technology Corporation LTspice XVII
						Backannotation: _in1 1 2
						Backannotation: _in2 1 2
						Backannotation: _in3 1 2
						Backannotation: _in4 1 2
						Backannotation: _in1 1 2
						Backannotation: _in2 1 2
						Backannotation: _in3 1 2
						Backannotation: _in4 1 2
						Backannotation: u1 1 2 99 50 45
						Variables:
							0	time	time
							1	V(+v)	voltage
							2	V(in)	voltage
							3	V(n001)	voltage
							4	V(out)	voltage
							5	V(-v)	voltage
							6	I(Rload)	device_current
							7	I(R2)	device_current
							8	I(R1)	device_current
							9	I(V2)	device_current
							10	I(Vin)	device_current
							11	I(V1)	device_current
							12	Ix(u1:1)	subckt_current
							13	Ix(u1:2)	subckt_current
							14	Ix(u1:99)	subckt_current
							15	Ix(u1:50)	subckt_current
							16	Ix(u1:45)	subckt_current
						Binary:
			`

	test_3_headers := `
						Title: * Z:\home\wine\ltspice\test_2.asc
						Date: Sat Jul 29 12:52:53 2023
						Plotname: Transient Analysis
						Flags: real forward
						No. Variables: 104
						No. Points:          841
						Offset:   0.0000000000000000e+000
						Command: Linear Technology Corporation LTspice XVII
						Variables:
							0	time	time
							1	V(n001)	voltage
							2	V(3)	voltage
							3	V(n008)	voltage
							4	V(2)	voltage
							5	V(n009)	voltage
							6	V(n011)	voltage
							7	V(n010)	voltage
							8	V(n013)	voltage
							9	V(n016)	voltage
							10	V(n018)	voltage
							11	V(n019)	voltage
							12	V(7)	voltage
							13	V(4)	voltage
							14	V(n002)	voltage
							15	V(n003)	voltage
							16	V(n014)	voltage
							17	V(n020)	voltage
							18	V(n004)	voltage
							19	V(n007)	voltage
							20	V(n015)	voltage
							21	V(n017)	voltage
							22	V(n005)	voltage
							23	V(n012)	voltage
							24	V(6)	voltage
							25	V(n006)	voltage
							26	Ic(Q19)	device_current
							27	Ib(Q19)	device_current
							28	Ie(Q19)	device_current
							29	Ic(Q11)	device_current
							30	Ib(Q11)	device_current
							31	Ie(Q11)	device_current
							32	Ic(Q10)	device_current
							33	Ib(Q10)	device_current
							34	Ie(Q10)	device_current
							35	Ic(Q9)	device_current
							36	Ib(Q9)	device_current
							37	Ie(Q9)	device_current
							38	Ic(Q4)	device_current
							39	Ib(Q4)	device_current
							40	Ie(Q4)	device_current
							41	Ic(Q6)	device_current
							42	Ib(Q6)	device_current
							43	Ie(Q6)	device_current
							44	Ic(Q5)	device_current
							45	Ib(Q5)	device_current
							46	Ie(Q5)	device_current
							47	Ic(Q20)	device_current
							48	Ib(Q20)	device_current
							49	Ie(Q20)	device_current
							50	Ic(Q18)	device_current
							51	Ib(Q18)	device_current
							52	Ie(Q18)	device_current
							53	Ic(Q17)	device_current
							54	Ib(Q17)	device_current
							55	Ie(Q17)	device_current
							56	Ic(Q16)	device_current
							57	Ib(Q16)	device_current
							58	Ie(Q16)	device_current
							59	Ic(Q15)	device_current
							60	Ib(Q15)	device_current
							61	Ie(Q15)	device_current
							62	Ic(Q14)	device_current
							63	Ib(Q14)	device_current
							64	Ie(Q14)	device_current
							65	Ic(Q13)	device_current
							66	Ib(Q13)	device_current
							67	Ie(Q13)	device_current
							68	Ic(Q12)	device_current
							69	Ib(Q12)	device_current
							70	Ie(Q12)	device_current
							71	Ic(Q3)	device_current
							72	Ib(Q3)	device_current
							73	Ie(Q3)	device_current
							74	Ic(Q8)	device_current
							75	Ib(Q8)	device_current
							76	Ie(Q8)	device_current
							77	Ic(Q7)	device_current
							78	Ib(Q7)	device_current
							79	Ie(Q7)	device_current
							80	Ic(Q2)	device_current
							81	Ib(Q2)	device_current
							82	Ie(Q2)	device_current
							83	Ic(Q1)	device_current
							84	Ib(Q1)	device_current
							85	Ie(Q1)	device_current
							86	I(C1)	device_current
							87	I(R14)	device_current
							88	I(R13)	device_current
							89	I(R12)	device_current
							90	I(R11)	device_current
							91	I(R10)	device_current
							92	I(R9)	device_current
							93	I(R8)	device_current
							94	I(R7)	device_current
							95	I(R6)	device_current
							96	I(R5)	device_current
							97	I(R4)	device_current
							98	I(R3)	device_current
							99	I(R2)	device_current
							100	I(R1)	device_current
							101	I(V3)	device_current
							102	I(V2)	device_current
							103	I(V1)	device_current
						Binary:
`
	tests := []struct {
		name     string
		input    string
		expected *MetaData
	}{
		{
			name:  "Test 1",
			input: test_1_headers,
			expected: &MetaData{
				Title:       "* Z:\\home\\wine\\ltspice\\Draft1.asc",
				Date:        time.Date(2023, 7, 25, 12, 15, 28, 0, time.UTC),
				SimType:     TransientAnalysis,
				Flags:       None | Forward | Stepped,
				NoVariables: 4,
				NoPoints:    142,
				Offset:      0,
				Command:     "Linear Technology Corporation LTspice XVII",
				Variables: []Variable{
					{
						order: 0,
						Name:  "time",
						Typ:   "time",
						size:  8,
					},
					{
						order: 1,
						Name:  "V(n001)",
						Typ:   "voltage",
						size:  4,
					},
					{
						order: 2,
						Name:  "I(R1)",
						Typ:   "device_current",
						size:  4,
					},
					{
						order: 3,
						Name:  "I(V1)",
						Typ:   "device_current",
						size:  4,
					},
				},
				BinaryOffset: 0,
			},
		},
		{
			name:  "Test 2",
			input: test_2_headers,
			expected: &MetaData{
				Title:       "* Z:\\home\\wine\\ltspice\\test.asc",
				Date:        time.Date(2023, 7, 29, 12, 38, 42, 0, time.UTC),
				SimType:     TransientAnalysis,
				Flags:       None | Forward,
				NoVariables: 17,
				NoPoints:    520,
				Offset:      0,
				Command:     "Linear Technology Corporation LTspice XVII",
				Variables: []Variable{
					{
						order: 0,
						Name:  "time",
						Typ:   "time",
						size:  8,
					},
					{
						order: 1,
						Name:  "V(+v)",
						Typ:   "voltage",
						size:  4,
					},
					{
						order: 2,
						Name:  "V(in)",
						Typ:   "voltage",
						size:  4,
					},
					{
						order: 3,
						Name:  "V(n001)",
						Typ:   "voltage",
						size:  4,
					},
					{
						order: 4,
						Name:  "V(out)",
						Typ:   "voltage",
						size:  4,
					},
					{
						order: 5,
						Name:  "V(-v)",
						Typ:   "voltage",
						size:  4,
					},
					{
						order: 6,
						Name:  "I(Rload)",
						Typ:   "device_current",
						size:  4,
					},
					{
						order: 7,
						Name:  "I(R2)",
						Typ:   "device_current",
						size:  4,
					},
					{
						order: 8,
						Name:  "I(R1)",
						Typ:   "device_current",
						size:  4,
					},
					{
						order: 9,
						Name:  "I(V2)",
						Typ:   "device_current",
						size:  4,
					},
					{
						order: 10,
						Name:  "I(Vin)",
						Typ:   "device_current",
						size:  4,
					},
					{
						order: 11,
						Name:  "I(V1)",
						Typ:   "device_current",
						size:  4,
					},
					{
						order: 12,
						Name:  "Ix(u1:1)",
						Typ:   "subckt_current",
						size:  4,
					},
					{
						order: 13,
						Name:  "Ix(u1:2)",
						Typ:   "subckt_current",
						size:  4,
					},
					{
						order: 14,
						Name:  "Ix(u1:99)",
						Typ:   "subckt_current",
						size:  4,
					},
					{
						order: 15,
						Name:  "Ix(u1:50)",
						Typ:   "subckt_current",
						size:  4,
					},
					{
						order: 16,
						Name:  "Ix(u1:45)",
						Typ:   "subckt_current",
						size:  4,
					},
				},
				BinaryOffset: 0,
			},
		},
		{
			name:  "Test 3",
			input: test_3_headers,
			expected: &MetaData{
				Title:       "* Z:\\home\\wine\\ltspice\\test_2.asc",
				Date:        time.Date(2023, 7, 29, 12, 52, 53, 0, time.UTC),
				SimType:     TransientAnalysis,
				Flags:       None | Forward,
				NoVariables: 104,
				NoPoints:    841,
				Offset:      0,
				Command:     "Linear Technology Corporation LTspice XVII",
				Variables: []Variable{
					{order: 0, Name: "time", Typ: "time", size: 8},
					{order: 1, Name: "V(n001)", Typ: "voltage", size: 4},
					{order: 2, Name: "V(3)", Typ: "voltage", size: 4},
					{order: 3, Name: "V(n008)", Typ: "voltage", size: 4},
					{order: 4, Name: "V(2)", Typ: "voltage", size: 4},
					{order: 5, Name: "V(n009)", Typ: "voltage", size: 4},
					{order: 6, Name: "V(n011)", Typ: "voltage", size: 4},
					{order: 7, Name: "V(n010)", Typ: "voltage", size: 4},
					{order: 8, Name: "V(n013)", Typ: "voltage", size: 4},
					{order: 9, Name: "V(n016)", Typ: "voltage", size: 4},
					{order: 10, Name: "V(n018)", Typ: "voltage", size: 4},
					{order: 11, Name: "V(n019)", Typ: "voltage", size: 4},
					{order: 12, Name: "V(7)", Typ: "voltage", size: 4},
					{order: 13, Name: "V(4)", Typ: "voltage", size: 4},
					{order: 14, Name: "V(n002)", Typ: "voltage", size: 4},
					{order: 15, Name: "V(n003)", Typ: "voltage", size: 4},
					{order: 16, Name: "V(n014)", Typ: "voltage", size: 4},
					{order: 17, Name: "V(n020)", Typ: "voltage", size: 4},
					{order: 18, Name: "V(n004)", Typ: "voltage", size: 4},
					{order: 19, Name: "V(n007)", Typ: "voltage", size: 4},
					{order: 20, Name: "V(n015)", Typ: "voltage", size: 4},
					{order: 21, Name: "V(n017)", Typ: "voltage", size: 4},
					{order: 22, Name: "V(n005)", Typ: "voltage", size: 4},
					{order: 23, Name: "V(n012)", Typ: "voltage", size: 4},
					{order: 24, Name: "V(6)", Typ: "voltage", size: 4},
					{order: 25, Name: "V(n006)", Typ: "voltage", size: 4},
					{order: 26, Name: "Ic(Q19)", Typ: "device_current", size: 4},
					{order: 27, Name: "Ib(Q19)", Typ: "device_current", size: 4},
					{order: 28, Name: "Ie(Q19)", Typ: "device_current", size: 4},
					{order: 29, Name: "Ic(Q11)", Typ: "device_current", size: 4},
					{order: 30, Name: "Ib(Q11)", Typ: "device_current", size: 4},
					{order: 31, Name: "Ie(Q11)", Typ: "device_current", size: 4},
					{order: 32, Name: "Ic(Q10)", Typ: "device_current", size: 4},
					{order: 33, Name: "Ib(Q10)", Typ: "device_current", size: 4},
					{order: 34, Name: "Ie(Q10)", Typ: "device_current", size: 4},
					{order: 35, Name: "Ic(Q9)", Typ: "device_current", size: 4},
					{order: 36, Name: "Ib(Q9)", Typ: "device_current", size: 4},
					{order: 37, Name: "Ie(Q9)", Typ: "device_current", size: 4},
					{order: 38, Name: "Ic(Q4)", Typ: "device_current", size: 4},
					{order: 39, Name: "Ib(Q4)", Typ: "device_current", size: 4},
					{order: 40, Name: "Ie(Q4)", Typ: "device_current", size: 4},
					{order: 41, Name: "Ic(Q6)", Typ: "device_current", size: 4},
					{order: 42, Name: "Ib(Q6)", Typ: "device_current", size: 4},
					{order: 43, Name: "Ie(Q6)", Typ: "device_current", size: 4},
					{order: 44, Name: "Ic(Q5)", Typ: "device_current", size: 4},
					{order: 45, Name: "Ib(Q5)", Typ: "device_current", size: 4},
					{order: 46, Name: "Ie(Q5)", Typ: "device_current", size: 4},
					{order: 47, Name: "Ic(Q20)", Typ: "device_current", size: 4},
					{order: 48, Name: "Ib(Q20)", Typ: "device_current", size: 4},
					{order: 49, Name: "Ie(Q20)", Typ: "device_current", size: 4},
					{order: 50, Name: "Ic(Q18)", Typ: "device_current", size: 4},
					{order: 51, Name: "Ib(Q18)", Typ: "device_current", size: 4},
					{order: 52, Name: "Ie(Q18)", Typ: "device_current", size: 4},
					{order: 53, Name: "Ic(Q17)", Typ: "device_current", size: 4},
					{order: 54, Name: "Ib(Q17)", Typ: "device_current", size: 4},
					{order: 55, Name: "Ie(Q17)", Typ: "device_current", size: 4},
					{order: 56, Name: "Ic(Q16)", Typ: "device_current", size: 4},
					{order: 57, Name: "Ib(Q16)", Typ: "device_current", size: 4},
					{order: 58, Name: "Ie(Q16)", Typ: "device_current", size: 4},
					{order: 59, Name: "Ic(Q15)", Typ: "device_current", size: 4},
					{order: 60, Name: "Ib(Q15)", Typ: "device_current", size: 4},
					{order: 61, Name: "Ie(Q15)", Typ: "device_current", size: 4},
					{order: 62, Name: "Ic(Q14)", Typ: "device_current", size: 4},
					{order: 63, Name: "Ib(Q14)", Typ: "device_current", size: 4},
					{order: 64, Name: "Ie(Q14)", Typ: "device_current", size: 4},
					{order: 65, Name: "Ic(Q13)", Typ: "device_current", size: 4},
					{order: 66, Name: "Ib(Q13)", Typ: "device_current", size: 4},
					{order: 67, Name: "Ie(Q13)", Typ: "device_current", size: 4},
					{order: 68, Name: "Ic(Q12)", Typ: "device_current", size: 4},
					{order: 69, Name: "Ib(Q12)", Typ: "device_current", size: 4},
					{order: 70, Name: "Ie(Q12)", Typ: "device_current", size: 4},
					{order: 71, Name: "Ic(Q3)", Typ: "device_current", size: 4},
					{order: 72, Name: "Ib(Q3)", Typ: "device_current", size: 4},
					{order: 73, Name: "Ie(Q3)", Typ: "device_current", size: 4},
					{order: 74, Name: "Ic(Q8)", Typ: "device_current", size: 4},
					{order: 75, Name: "Ib(Q8)", Typ: "device_current", size: 4},
					{order: 76, Name: "Ie(Q8)", Typ: "device_current", size: 4},
					{order: 77, Name: "Ic(Q7)", Typ: "device_current", size: 4},
					{order: 78, Name: "Ib(Q7)", Typ: "device_current", size: 4},
					{order: 79, Name: "Ie(Q7)", Typ: "device_current", size: 4},
					{order: 80, Name: "Ic(Q2)", Typ: "device_current", size: 4},
					{order: 81, Name: "Ib(Q2)", Typ: "device_current", size: 4},
					{order: 82, Name: "Ie(Q2)", Typ: "device_current", size: 4},
					{order: 83, Name: "Ic(Q1)", Typ: "device_current", size: 4},
					{order: 84, Name: "Ib(Q1)", Typ: "device_current", size: 4},
					{order: 85, Name: "Ie(Q1)", Typ: "device_current", size: 4},
					{order: 86, Name: "I(C1)", Typ: "device_current", size: 4},
					{order: 87, Name: "I(R14)", Typ: "device_current", size: 4},
					{order: 88, Name: "I(R13)", Typ: "device_current", size: 4},
					{order: 89, Name: "I(R12)", Typ: "device_current", size: 4},
					{order: 90, Name: "I(R11)", Typ: "device_current", size: 4},
					{order: 91, Name: "I(R10)", Typ: "device_current", size: 4},
					{order: 92, Name: "I(R9)", Typ: "device_current", size: 4},
					{order: 93, Name: "I(R8)", Typ: "device_current", size: 4},
					{order: 94, Name: "I(R7)", Typ: "device_current", size: 4},
					{order: 95, Name: "I(R6)", Typ: "device_current", size: 4},
					{order: 96, Name: "I(R5)", Typ: "device_current", size: 4},
					{order: 97, Name: "I(R4)", Typ: "device_current", size: 4},
					{order: 98, Name: "I(R3)", Typ: "device_current", size: 4},
					{order: 99, Name: "I(R2)", Typ: "device_current", size: 4},
					{order: 100, Name: "I(R1)", Typ: "device_current", size: 4},
					{order: 101, Name: "I(V3)", Typ: "device_current", size: 4},
					{order: 102, Name: "I(V2)", Typ: "device_current", size: 4},
					{order: 103, Name: "I(V1)", Typ: "device_current", size: 4},
				},
				BinaryOffset: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			EncodedHeaders := utf16.Encode([]rune(tt.input))
			buff := new(bytes.Buffer)
			for _, c := range EncodedHeaders {
				err := binary.Write(buff, binary.LittleEndian, c)
				if err != nil {
					t.Fatal(err)
				}
			}

			reader := bytes.NewReader(buff.Bytes())

			got, err := parseHeaders(reader)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expected.Title, got.Title, "Title mismatch")
			assert.Equal(t, tt.expected.Flags, got.Flags, "Flags mismatch")
			assert.Equal(t, tt.expected.NoPoints, got.NoPoints, "No. Of Points mismatch")
			assert.Equal(t, tt.expected.NoVariables, got.NoVariables, "No. of variables mismatch")
			assert.Equal(t, tt.expected.Command, got.Command, "command mismatch")
			assert.Equal(t, tt.expected.Offset, got.Offset, "Offset mismatch")
			assert.Equal(t, tt.expected.SimType, got.SimType, "Simulation type mismatch")
			assert.ElementsMatch(t, tt.expected.Variables, got.Variables, "Variables mismatch")
			assert.True(t, tt.expected.Date.Equal(got.Date), "Date mismatch")

		})
	}
}
