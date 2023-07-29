package main

import (
	"fmt"

	"github.com/theadell/ltspice-go/ltspice"
)

func main() {
	sim, err := ltspice.Parse("dc-sweep.raw")
	if err != nil {
		panic(err)
	}
	for k, v := range sim.Data {
		fmt.Printf("-%-10s has length of %-3d\n", k, len(v))
	}

}
