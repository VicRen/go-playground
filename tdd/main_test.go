package main

import (
	"reflect"
	"testing"
)

func TestPrimeFactor(t *testing.T) {
	tt := []struct {
		name string
		parm int
		want []int
	}{
		{
			"prime factor of one",
			1,
			make([]int, 0),
		},
		{
			"prime factor of two",
			2,
			[]int{2},
		},
		{
			"prime factor of three",
			3,
			[]int{3},
		},
		{
			"prime factor of four",
			4,
			[]int{2, 2},
		},
		{
			"prime factor of five",
			5,
			[]int{5},
		},
		{
			"prime factor of six",
			6,
			[]int{2, 3},
		},
		{
			"prime factor of seven",
			7,
			[]int{7},
		},
		{
			"prime factor of eight",
			8,
			[]int{2, 2, 2},
		},
		{
			"prime factor of nine",
			9,
			[]int{3, 3},
		},
		{
			"prime factor of nine",
			2 * 2 * 2 * 3 * 3 * 3 * 5 * 7 * 11 * 13 * 17 * 23,
			[]int{2, 2, 2, 3, 3, 3, 5, 7, 11, 13, 17, 23},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := PrimeFactorOf(tc.parm)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("PrimeFactor of %d should be %v, got %v", tc.parm, tc.want, got)
			}
		})
	}
}

func PrimeFactorOf(x int) []int {
	ret := make([]int, 0)
	divider := 2
	for divider < x {
		for x%divider == 0 {
			ret = append(ret, divider)
			x = x / divider
		}
		divider++
	}
	if x > 1 {
		ret = append(ret, x)
	}
	return ret
}
