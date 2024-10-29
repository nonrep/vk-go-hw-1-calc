package main

import (
	"fmt"
	"os"

	"github.com/nonrep/go-homework-1-calc/calc"
)

func main() {
	formula := os.Args[1]
	result, err := calc.Calc(formula)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}
