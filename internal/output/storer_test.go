package output_test

import (
	"github.com/YReshetko/go-annotation/internal/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

const tempDir = "./tmp"
const expected = `// Code generated by TEST annotation processor. DO NOT EDIT.
// versions:
//		go-annotation: 0.0.2-alpha
//		TEST: 1.0.0
package fixtures
`

func TestTextStorer_Success(t *testing.T) {
	defer clean(t)

	file := filepath.Join(tempDir, "dummy.go")
	err := output.Store(file, []byte("package fixtures\n"), map[string]string{"annotation_name": "TEST", "annotation_ver": "1.0.0"})
	require.NoError(t, err)

	data, err := os.ReadFile(file)
	require.NoError(t, err)

	assert.Equal(t, expected, string(data))
}

func TestASTStorer_Success(t *testing.T) {
	defer clean(t)
	file := filepath.Join(tempDir, "dummy.go")
	f := &ast.File{
		Package: token.Pos(1),
		Name: &ast.Ident{
			NamePos: 9,
			Name:    "fixtures",
		},
	}

	err := output.StoreAST(file, f, map[string]string{"annotation_name": "TEST", "annotation_ver": "1.0.0"})
	require.NoError(t, err)

	data, err := os.ReadFile(file)
	require.NoError(t, err)

	assert.Equal(t, expected, string(data))
}

func clean(t *testing.T) {
	require.NoError(t, os.RemoveAll(tempDir))
}
