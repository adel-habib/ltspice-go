# ltspice-go
A Go-based parser for LTSpice raw binary data.  
## TODOs 

- [ ] Core Features
    - [x] Parse LTSpice raw file metadata
    - [x] Parse LTSpice raw file binary data
    - [ ] Write LTSpice raw file metadata
    - [ ] Write LTSpice raw file binary data

- [ ] Additional Features
    - [x] Handle complex values in binary data
    - [ ] Handle fast access data structure in binary data
    - [ ] Handle stepped simulations (extract stepping information from .log files)

- [ ] Data Analysis and Utilities
    - [ ] Provide functions to manipulate parsed simulation data
        - [ ] Filter by variable
        - [ ] Filter by time range
    - [ ] Provide functions to export data to other formats (CSV, JSON, etc.)
    - [ ] Provide functions to generate plots

- [x] Simulations supported
    - [x] Operation Point
    - [x] DC transfer characteristic
    - [x] AC Analysis
    - [x] Transient Analysis
    - [x] Transfer Function
    - [x] Noise Spectral Density - (V/Hz½ or A/Hz½)

- [ ] Documentation and Examples
    - [ ] Write comprehensive documentation
    - [ ] Provide usage examples

- [ ] Testing and Validation
    - [x] Unit tests
    - [x] Test against a variety of LTSpice raw files
    - [ ] Validate correct parsing of binary data
    - [ ] Validate correct writing of binary data
