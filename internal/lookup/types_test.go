package lookup_test

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/YReshetko/go-annotation/internal/lookup"
	"github.com/YReshetko/go-annotation/internal/module"
)

func TestFindType(t *testing.T) {
	m, err := module.Load("./fixtures")
	require.NoError(t, err)

	node, parents, _, filePath, err := lookup.FindTypeByImport(m, "github.com/YReshetko/go-annotation/internal/lookup/fixtures/dashed-package", "Exported2")
	require.NoError(t, err)
	assert.Contains(t, filePath, "src/github.com/YReshetko/go-annotation/internal/lookup/fixtures/dashed-package/example2.go")
	assert.Len(t, parents, 3)

	ts, ok := node.(*ast.TypeSpec)
	require.True(t, ok)

	require.Equal(t, "Exported2", ts.Name.Name)
	_, ok = ts.Type.(*ast.StructType)
	require.True(t, ok)
}

func TestFindType_NotFound(t *testing.T) {
	m, err := module.Load("./fixtures")
	require.NoError(t, err)

	node, parents, _, _, err := lookup.FindTypeByImport(m, "github.com/YReshetko/go-annotation/internal/lookup/fixtures/dashed-package", "Exported3")
	require.Error(t, err)
	assert.Nil(t, node)
	assert.Empty(t, parents)
}

func TestFindTypeInDir(t *testing.T) {
	m, err := module.Load("./fixtures")
	require.NoError(t, err)

	node, parents, _, filePath, err := lookup.FindTypeInDir(m, "dashed-package", "Exported2")
	require.NoError(t, err)
	assert.Contains(t, filePath, "src/github.com/YReshetko/go-annotation/internal/lookup/fixtures/dashed-package/example2.go")
	assert.Len(t, parents, 3)

	ts, ok := node.(*ast.TypeSpec)
	require.True(t, ok)

	require.Equal(t, "Exported2", ts.Name.Name)
	_, ok = ts.Type.(*ast.StructType)
	require.True(t, ok)
}

func TestFindTypeInDir_NotFound(t *testing.T) {
	m, err := module.Load("./fixtures")
	require.NoError(t, err)

	node, parents, _, _, err := lookup.FindTypeInDir(m, "dashed-package", "Exported3")
	require.Error(t, err)
	assert.Nil(t, node)
	assert.Empty(t, parents)
}
