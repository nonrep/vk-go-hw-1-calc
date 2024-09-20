package main

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func TestOK(t *testing.T) {
	out := bytes.NewBuffer(nil)
	var tests = []struct {
		formula string
		expect  int
	}{
		{"1", 1},
		{"1 + 2", 3},
		{"1 - 13", -12},
		{"2 * 2", 4},
		{"10 / 2", 5},
		{"1 + 2 * 3", 7},
		{" (     (  (    10 )  )  ) ", 10},
		{"5  *  4 * 3 * ((2 * 1) + 7) * 8", 4320},
		{"15/(7-(1 + 1))*3", 9},
		{"    ( ( (1 + 2)  )  )  * 10 ", 30},
		{"15/(7-(1+1))*3-(2+(1+1))*15/(7-(200+1))*3-(2+(1+1))*(15/(7-(1+1))*3-(2+(1+1))+15/(7-(1+1))*3-(2+(1+1)))", -31},
	}

	for _, test := range tests {
		err := calc(test.formula, out)
		if err != nil {
			t.Errorf("Testing OK failed: %s", err)
		}
		data := strings.TrimSpace(out.String())
		result, err := strconv.Atoi(data)
		if err != nil {
			t.Errorf("Testing OK failed: %s", err)
		}

		if result != test.expect {
			t.Error("Testing OK failed, result not match")
		}
		out.Reset()
	}
}

func TestFail(t *testing.T) {
	out := bytes.NewBuffer(nil)
	var tests = []struct {
		formula string
	}{
		{"1 2 3"},
		{"1+"},
		{"1/"},
		{"1 2 /"},
		{"1  + 2 /"},
		{"1 / 0"},
		{"* 1 2"},
		{")(1)"},
		{"()"},
		{"12)"},
		{"(12"},
		{""},
		{"abc"},
		{"a+c"},
		{""},
	}

	for _, test := range tests {
		err := calc(test.formula, out)
		if err == nil {
			t.Errorf("Test FAIL failed: expected error")
		}
	}
}
