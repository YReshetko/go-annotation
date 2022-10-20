package stream_test

import (
	"github.com/YReshetko/go-annotation/internal/utils/stream"
	"testing"
)

func genArr(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i
	}
	return arr
}

var testCases = []struct {
	name  string
	value []int
}{
	{"100", genArr(100)},
	{"1000", genArr(1000)},
	{"10000", genArr(10000)},
	{"100000", genArr(100000)},
	{"1000000", genArr(1000000)},
	/*	{"10000000", genArr(10000000)},
		{"100000000", genArr(100000000)},
		{"1000000000", genArr(1000000000)},*/
}

func BenchmarkImperative(b *testing.B) {
	for _, v := range testCases {
		b.Run(v.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				var result []int
				for _, v := range v.value {
					if v%3 == 0 {
						result = append(result, v*v)
					}
				}
				_ = result
			}
		})
	}
}

func BenchmarkFunctional(b *testing.B) {
	for _, v := range testCases {
		b.Run(v.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = stream.Map(stream.OfSlice(v.value).
					Filter(func(n int) bool {
						return n%3 == 0
					}), func(n int) int {
					return n * n
				}).ToSlice()
			}
		})
	}
}
