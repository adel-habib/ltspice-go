# ltspice-go
Ltspice parser for Go

## TODOs 

- [ ] Core Features
    - [x] Parse LTSpice raw file metadata
    - [x] Parse LTSpice raw file binary data
    - [ ] Write LTSpice raw file metadata
    - [ ] Write LTSpice raw file binary data

- [ ] Additional Features
    - [ ] Handle complex values in binary data
    - [ ] Handle fast access data structure in binary data
    - [ ] Handle stepped simulations (extract stepping information from .log files)

- [ ] Data Analysis and Utilities
    - [ ] Provide functions to manipulate parsed simulation data
        - [ ] Filter by variable
        - [ ] Filter by time range
    - [ ] Provide functions to export data to other formats (CSV, JSON, etc.)
    - [ ] Provide functions to generate plots

- [ ] Simulation Specific Features
    - [ ] Operation Point
        - [ ] TODO: Add specific features related to Operation Point simulation
    - [ ] DC transfer characteristic
        - [ ] TODO: Add specific features related to DC transfer characteristic simulation
    - [ ] AC Analysis
        - [ ] TODO: Add specific features related to AC Analysis simulation
    - [ ] Transient Analysis
        - [x] Parse binary data for unstepped simulations
        - [ ] Parse binary data for stepped simulations
    - [ ] Noise Spectral Density - (V/Hz½ or A/Hz½)
        - [ ] TODO: Add specific features related to Noise Spectral Density simulation
    - [ ] Transfer Function
        - [ ] TODO: Add specific features related to Transfer Function simulation

- [ ] Documentation and Examples
    - [ ] Write comprehensive documentation
    - [ ] Provide usage examples

- [ ] Testing and Validation
    - [x] Unit tests
    - [ ] Test against a variety of LTSpice raw files
    - [ ] Validate correct parsing and writing of binary data
