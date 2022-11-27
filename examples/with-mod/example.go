package with_mod

import (
	"github.com/eapache/go-resiliency/breaker"
	_ "github.com/eapache/go-resiliency/breaker"
	"github.com/eapache/queue"
	_ "github.com/eapache/queue"
)

// @Mapper
type SomeInterface interface {
	mapper1(queue.Queue) queue.Queue
	mapper(breaker.Breaker) breaker.Breaker
}
