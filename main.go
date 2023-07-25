package main

import (
	"fmt"

	"github.com/theadell/ltspice-go/ltspice"
)

func main() {
	m, err := ltspice.Parse("iter.raw")
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
}
