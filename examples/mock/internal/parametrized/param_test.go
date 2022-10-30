package parametrized_test

import (
	"github.com/YReshetko/go-annotation/examples/mock/internal/parametrized"
	"github.com/YReshetko/go-annotation/examples/mock/internal/parametrized/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcess(t *testing.T) {
	m := &mocks.GenericInterfaceMock{}
	ch := make(chan int, 13)
	for i := 0; i < 13; i++ {
		ch <- i
	}
	close(ch)
	m.ProcessReturns(ch, 0)
	ok := parametrized.Process[int, float32](m, 1, []float32{10.0})
	assert.True(t, ok)
}
