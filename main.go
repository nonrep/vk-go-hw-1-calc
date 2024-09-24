package main

import (
	"fmt"
	"os"
)

func main() {
	formula := os.Args[1]
	result, err := Calc(formula)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}
