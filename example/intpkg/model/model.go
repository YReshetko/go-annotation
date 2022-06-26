package model

import (
	"github.com/YReshetko/go-annotation/annotations/rest"
	"github.com/YReshetko/interview-tasks/pkg/heap-pattern"
)

type User struct {
	r rest.Rest
	h heap_pattern.Heap
}

type ExternalFunction func() bool
