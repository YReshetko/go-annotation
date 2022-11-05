// Invalid example but that would be nice to explore the case deeply

package builder

import (
	"math/rand"
	"net/http"
)

type Builder[V Builder[V, T], T any] interface {
	Field1(v http.Client) V
	Field2(v int) V
	Field3(v string) V
	Build() T
}

// @Constructor
type Factory[V Builder[V, T], T any] struct {
	registry map[string]func() Builder[V, T] // @Init
}

func (f *Factory[V, T]) Register(key string, fn func() Builder[V, T]) {
	f.registry[key] = fn
}

func (f *Factory[V, T]) New(key string) T {
	var v T
	if fn, ok := f.registry[key]; ok {
		builder := fn()
		return builder.Field2(rand.Int()).Field3("set value").Build()
	}
	return v
}
