package module

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTree_Success(t *testing.T) {
	testTree := newTree()
	toTest := []struct {
		s selector
		m *Module
	}{
		{
			s: selector{pkg: "github.com/my/root_project"},
			m: &Module{},
		},
		{
			s: selector{pkg: "github.com/my/root_project", next: &selector{pkg: "golang.org/my/child_1"}},
			m: &Module{},
		},
		{
			s: selector{pkg: "github.com/my/root_project", next: &selector{pkg: "golang.org/my/child_2"}},
			m: &Module{},
		},
		{
			s: selector{pkg: "github.com/my/root_project", next: &selector{pkg: "golang.org/my/child_2", next: &selector{pkg: "github.com/my/child_2_1"}}},
			m: &Module{},
		},
		{
			s: selector{pkg: "github.com/my/root_project", next: &selector{pkg: "golang.org/my/child_2", next: &selector{pkg: "github.com/my/child_2_2"}}},
			m: &Module{},
		},
		{
			s: selector{pkg: "github.com/my/root_project", next: &selector{pkg: "golang.org/my/child_2", next: &selector{pkg: "github.com/my/child_2_3"}}},
			m: &Module{},
		},
	}

	for _, tt := range toTest {
		err := testTree.put(tt.s, tt.m)
		require.NoError(t, err)
	}

	for _, tt := range toTest {
		m, ok := testTree.get(tt.s)
		require.True(t, ok)
		assert.True(t, tt.m == m)
	}
}

func TestTree_FailOnPut(t *testing.T) {
	testTree := newTree()
	s := selector{pkg: "github.com/my/root_project", next: &selector{pkg: "golang.org/my/child_1"}}

	err := testTree.put(s, &Module{})
	require.Error(t, err)
	assert.EqualError(t, err, "unable to create modules tree node as no existing module for github.com/my/root_project")
}

func TestTree_FailOnReplace(t *testing.T) {
	testTree := newTree()
	s1 := selector{pkg: "github.com/my/root_project"}
	s2 := selector{pkg: "github.com/my/root_project", next: &selector{pkg: "golang.org/my/child_1"}}
	s3 := selector{pkg: "github.com/my/root_project", next: &selector{pkg: "golang.org/my/child_1"}}

	err := testTree.put(s1, &Module{})
	require.NoError(t, err)
	err = testTree.put(s2, &Module{})
	require.NoError(t, err)
	err = testTree.put(s3, &Module{})
	require.Error(t, err)
	assert.EqualError(t, err, "unable to replace module in tree for pkg golang.org/my/child_1")
}

func TestTree_NotFound(t *testing.T) {
	testTree := newTree()
	s1 := selector{pkg: "github.com/my/root_project"}
	s2 := selector{pkg: "github.com/my/root_project", next: &selector{pkg: "golang.org/my/child_1"}}
	s3 := selector{pkg: "github.com/my/root_project", next: &selector{pkg: "golang.org/my/child_2"}}

	err := testTree.put(s1, &Module{})
	require.NoError(t, err)
	err = testTree.put(s2, &Module{})
	require.NoError(t, err)

	m, ok := testTree.get(s3)
	require.False(t, ok)
	require.Nil(t, m)
}
