package main

import (
	"github.com/theadell/ltspice-go/ltspice"
)

func main() {
	_, err := ltspice.Parse("rc.raw")
	if err != nil {
		panic(err)
	}
	//fmt.Println(m)
}
