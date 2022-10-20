package ast_test

import (
	ast2 "go/ast"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/YReshetko/go-annotation/internal/ast"
)

func TestNodeSearch(t *testing.T) {
	f, err := ast.LoadFileAST("./fixtures/test_file.go")
	require.NoError(t, err)

	toTest := []struct {
		name string
		node ast2.Node
	}{
		{name: "SingleConst", node: &ast2.ValueSpec{}},
		{name: "GroupConst1", node: &ast2.ValueSpec{}},
		{name: "GroupConst2", node: &ast2.ValueSpec{}},
		{name: "SingleVar", node: &ast2.ValueSpec{}},
		{name: "GroupVar1", node: &ast2.ValueSpec{}},
		{name: "GroupVar2", node: &ast2.ValueSpec{}},
		{name: "GroupSeveralVars1", node: &ast2.ValueSpec{}},
		{name: "GroupSeveralVars2", node: &ast2.ValueSpec{}},

		{name: "SingleInterface", node: &ast2.TypeSpec{}},
		{name: "SingleStruct", node: &ast2.TypeSpec{}},
		{name: "GroupInterface1", node: &ast2.TypeSpec{}},
		{name: "GroupStruct1", node: &ast2.TypeSpec{}},
		{name: "GroupInterface2", node: &ast2.TypeSpec{}},
		{name: "GroupStruct2", node: &ast2.TypeSpec{}},

		{name: "SomeFunction", node: &ast2.FuncDecl{}},
		{name: "SomeMethod", node: &ast2.FuncDecl{}},
	}

	for _, s := range toTest {
		t.Run(s.name, func(t *testing.T) {
			verifyType(t, s.name, f, s.node)
		})
	}
}

func TestNodeSearch_NotFound(t *testing.T) {
	f, err := ast.LoadFileAST("./fixtures/test_file.go")
	require.NoError(t, err)

	toTest := []string{"SingleVariable", "SingleConstant", "InternalFunction"}

	for _, s := range toTest {
		t.Run(s, func(t *testing.T) {
			a, ok := ast.FindTopNodeByName(f, s)
			require.False(t, ok)
			assert.Nil(t, a)
		})
	}
}

func TestComment(t *testing.T) {
	f, err := ast.LoadFileAST("./fixtures/test_file.go")
	require.NoError(t, err)

	toTest := []struct {
		name    string
		comment string
	}{
		{
			name:    "SingleInterface",
			comment: "SingleInterface single line comment",
		},
		{
			name: "SingleStruct",
			comment: `multiline comment
new line`,
		},
		{
			name: "GroupStruct1",
			comment: `Several single line comments
Several single line comments`,
		},
		{
			name:    "GroupInterface2",
			comment: ``,
		},
		{
			name:    "GroupInterface2",
			comment: ``,
		},
		{
			name:    "SingleConst",
			comment: `Single line comment on constant`,
		},
		{
			name: "SeveralVars1",
			comment: `Multiline comment on
Variables`,
		},
		{
			name: "SeveralVars2",
			comment: `Multiline comment on
Variables`,
		},
		{
			name: "GroupConst1",
			comment: `		Multiline comment GroupConst1
		Variables`,
		},
		{
			name:    "GroupConst2",
			comment: `Single line comment on GroupConst2`,
		},
		{
			name:    "GroupVar1",
			comment: `Single line comment on GroupVar1`,
		},
		{
			name: "GroupSeveralVars1",
			comment: `	   Multiline comment GroupSeveralVars1 and GroupSeveralVars2
	   Variables`,
		},
		{
			name: "GroupSeveralVars2",
			comment: `	   Multiline comment GroupSeveralVars1 and GroupSeveralVars2
	   Variables`,
		},
		{
			name:    "SomeFunction",
			comment: `Single line comment on SomeFunction`,
		},
		{
			name: "SomeMethod",
			comment: `Multiline comment on SomeMethod
Method`,
		},
	}

	for _, s := range toTest {
		t.Run(s.name, func(t *testing.T) {
			a, ok := ast.FindTopNodeByName(f, s.name)
			require.True(t, ok)

			comment, ok := ast.Comment(a)
			require.True(t, ok)
			assert.Equal(t, s.comment, comment)
		})
	}
}

func TestNonTopLevelNodeComment(t *testing.T) {
	f, err := ast.LoadFileAST("./fixtures/test_file.go")
	require.NoError(t, err)

	toTest := []struct {
		originNodeName string
		nodeName       string
		comment        string
	}{
		{
			originNodeName: "GroupInterface2",
			nodeName:       "GroupConst2",
			comment:        "Non top lavel node comment",
		},
		{
			originNodeName: "GroupStruct2",
			nodeName:       "GroupVar2",
			comment: `			Multiline comment on
			GroupStruct2.GroupVar2`,
		},
	}

	for _, s := range toTest {
		t.Run(s.originNodeName+"."+s.nodeName, func(t *testing.T) {
			a, ok := ast.FindTopNodeByName(f, s.originNodeName)
			require.True(t, ok)

			found := false
			ast2.Inspect(a, func(node ast2.Node) bool {
				switch nt := node.(type) {
				case *ast2.Field:
					if len(nt.Names) == 0 {
						break
					}
					if nt.Names[0].Name == s.nodeName {
						found = true
						comment, ok := ast.Comment(nt)
						require.True(t, ok)
						assert.Equal(t, s.comment, comment)
					}
				}
				return true
			})
			assert.True(t, found)
		})
	}

}

func TestWalk_FindInterfaceAndStructureFields(t *testing.T) {
	f, err := ast.LoadFileAST("./fixtures/test_file.go")
	require.NoError(t, err)

	foundInterfaceField := false
	foundStructureField := false
	foundSingleInterface := false
	foundSingleStructure := false
	ast.Walk(f, func(node ast2.Node) bool {
		switch nt := node.(type) {
		case *ast2.Field:
			if len(nt.Names) != 1 {
				return true
			}
			if nt.Names[0].Name == "GroupConst2" && nt.Doc != nil {
				foundInterfaceField = true
				assert.Equal(t, "Non top lavel node comment", strings.TrimRight(nt.Doc.Text(), "\n"))
			}
			if nt.Names[0].Name == "GroupVar2" && nt.Doc != nil {
				foundStructureField = true
				assert.Equal(t, `			Multiline comment on
			GroupStruct2.GroupVar2`, strings.TrimRight(nt.Doc.Text(), "\n"))
			}
		case *ast2.TypeSpec:
			if nt.Name == nil {
				return true
			}
			if nt.Name.Name == "SingleInterface" {
				foundSingleInterface = true
				assert.Equal(t, "SingleInterface single line comment", strings.TrimRight(nt.Doc.Text(), "\n"))
			}
			if nt.Name.Name == "SingleStruct" {
				foundSingleStructure = true
				assert.Equal(t, "multiline comment\nnew line", strings.TrimRight(nt.Doc.Text(), "\n"))
			}

		}
		return true
	})

	assert.True(t, foundInterfaceField)
	assert.True(t, foundStructureField)
	assert.True(t, foundSingleInterface)
	assert.True(t, foundSingleStructure)
}

func verifyType(t *testing.T, name string, f *ast2.File, n ast2.Node) {
	a, ok := ast.FindTopNodeByName(f, name)
	require.True(t, ok)
	assert.Equal(t, reflect.TypeOf(n), reflect.TypeOf(a))
}
