package module_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/YReshetko/go-annotation/internal/module"
)

const root = "../../"

func TestLoad_Success(t *testing.T) {
	m, err := module.Load(root)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(m.Files()), 10)

	absRoot, err := filepath.Abs(root)
	require.NoError(t, err)

	assert.Equal(t, absRoot, m.Root())

	assert.Contains(t, m.Files(), "internal/module/files.go")
	assert.Contains(t, m.Files(), "internal/module/interface.go")
	assert.Contains(t, m.Files(), "internal/module/interface_test.go")
	assert.Contains(t, m.Files(), "internal/module/lookup.go")
	assert.Contains(t, m.Files(), "internal/module/module.go")
}

func TestSubmodule_Success(t *testing.T) {
	m, err := module.Load(root)
	require.NoError(t, err)

	m, err = module.Find(m, "github.com/stretchr/testify/assert")
	require.NoError(t, err)
	assert.Contains(t, m.Root(), "pkg/mod/github.com/stretchr/testify@v1.8.0")

	assert.Contains(t, m.Files(), "assert/assertion_compare.go")
	assert.Contains(t, m.Files(), "assert/assertion_compare_can_convert.go")
	assert.Contains(t, m.Files(), "mock/mock.go")
	assert.Contains(t, m.Files(), "require/require.go")
	assert.Contains(t, m.Files(), "package_test.go")

	m, err = module.Find(m, "gopkg.in/yaml.v3")
	require.NoError(t, err)
	assert.Contains(t, m.Root(), "pkg/mod/gopkg.in/yaml.v3@v3.0.1")

	assert.Contains(t, m.Files(), "decode.go")
	assert.Contains(t, m.Files(), "encode.go")

}

func TestSubmodule_Fail_NonDependency(t *testing.T) {
	m, err := module.Load(root)
	require.NoError(t, err)

	m, err = module.Find(m, "github.com/stretchr/testify/assert")
	require.NoError(t, err)
	assert.Contains(t, m.Root(), "pkg/mod/github.com/stretchr/testify@v1.8.0")

	assert.Contains(t, m.Files(), "assert/assertion_compare.go")
	assert.Contains(t, m.Files(), "assert/assertion_compare_can_convert.go")
	assert.Contains(t, m.Files(), "mock/mock.go")
	assert.Contains(t, m.Files(), "require/require.go")
	assert.Contains(t, m.Files(), "package_test.go")

	m, err = module.Find(m, "gopkg.in/yaml.v3")
	require.NoError(t, err)
	assert.Contains(t, m.Root(), "pkg/mod/gopkg.in/yaml.v3@v3.0.1")

	assert.Contains(t, m.Files(), "decode.go")
	assert.Contains(t, m.Files(), "encode.go")

}
