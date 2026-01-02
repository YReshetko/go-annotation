package environment_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/YReshetko/go-annotation/internal/environment"
)

func TestGoPath(t *testing.T) {
	// TODO Make it independent on local environment
	assert.Equal(t, "/home/yury/go/go1.25.5", environment.GoPath())
}

func TestModPath(t *testing.T) {
	// TODO Make it independent on local environment
	assert.Equal(t, "/home/yury/go/go1.25.5/pkg/mod", environment.ModPath())
}
