package main

import (
	"fmt"
	"testing"
)

func TestOK(t *testing.T) {
	var tests = []struct {
		formula string
		expect  float64
	}{
		{"1", 1},
		{"1.0", 1},
		{"1 + 2", 3},
		{"1.0 + 2", 3},
		{"1 - 13", -12},
		{"2 * 2", 4},
		{"10 / 2", 5},
		{"1 + 2 * 3", 7},
		{"1.3 + 2.1 * 3.5", 8.65},
		{" (     (  (    10 )  )  ) ", 10},
		{"5  *  4 * 3 * ((2 * 1) + 7) * 8", 4320},
		{"15/(7-(1 + 1))*3", 9},
		{"    ( ( (1 + 2)  )  )  * 10 ", 30},
		{"-(-11-(1*20/2)-11/2*3)", 37.5},
		{"  - (  -  11  -  ( 1  *  20 /   2  ) -11 /  2 *3 )  ", 37.5},
		{"15/(7-(1+1))*3-(2+(1+1))*15/(7-(200+1))*3-(2+(1+1))*(15/(7-(1+1))*3-(2+(1+1))+15/(7-(1+1))*3-(2+(1+1)))", -30.072164948453608},
	}

	for _, test := range tests {
		result, err := Calc(test.formula)
		if err != nil {
			t.Errorf("Testing OK failed: %s", err)
		}

		if result != test.expect {
			t.Error("Testing OK failed, result not match")
		}
	}
}

func TestFail(t *testing.T) {
	var tests = []struct {
		formula string
	}{
		{"1 2 3"},
		{"1. 2 3"},
		{"1.0 2.3 3.5"},
		{"1+"},
		{"1.2+"},
		{"1/"},
		{"1 2 /"},
		{"1  + 2 /"},
		{"1 / 0"},
		{"1 / 0.0"},
		{"* 1 2"},
		{")(1)"},
		{"()"},
		{"12)"},
		{"(12"},
		{""},
		{"abc"},
		{"a+c"},
		{"."},
		{""},
	}

	for _, test := range tests {
		_, err := Calc(test.formula)
		if err == nil {
			fmt.Println(test.formula)
			t.Errorf("Test FAIL failed: expected error")
		}
	}
}
