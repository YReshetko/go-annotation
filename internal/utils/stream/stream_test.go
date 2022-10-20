package stream_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/YReshetko/go-annotation/internal/utils/stream"
)

func TestStream(t *testing.T) {
	type a struct{ v int }
	type b struct{ v float64 }

	aToInt := func(t a) int { return t.v }

	intToB := func(t int) b { return b{v: float64(t)} }
	bToFloat := func(t b) float64 { return t.v }

	arr := []a{{v: 1}, {v: 3}, {v: 6}, {v: 9}, {v: 12}, {v: 15}, {v: 18}, {v: 21}, {v: 24}, {v: 24}, {v: 24}}
	result :=
		Map(
			FlatMap(
				Map(
					Map(
						Map(
							OfSlice(arr),
							aToInt,
						).Filter(mod5),
						intToB,
					),
					bToFloat,
				),
				splitMod5Float64,
			),
			double,
		).ToSlice()

	require.Len(t, result, 3)
	assert.Equal(t, float64(30), result[0])
	assert.Equal(t, float64(15), result[1])
	assert.Equal(t, float64(10), result[2])

}

func double(f float64) float64 {
	return 2 * f
}

func splitMod5Float64(t float64) []float64 {
	var o []float64
	for i := 0; i < int(t)/5; i++ {
		o = append(o, t/float64(i+1))
	}
	return o
}

func mod5(v int) bool {
	return v%5 == 0
}
