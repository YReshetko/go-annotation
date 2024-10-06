package imports_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/YReshetko/go-annotation/internal/utils/imports"
)

func TestOf(t *testing.T) {
	expectedPath := imports.Path("some/path/to/test/")
	actualPath := imports.Of("some/path/to/test/")
	assert.Equal(t, expectedPath, actualPath)
}

func TestPath_String(t *testing.T) {
	testCases := []struct {
		original string
		expected string
	}{
		{"some/path/to/test/", "some/path/to/test/"},
		{"some\\path\\to\\test", "some\\path\\to\\test"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.original, func(t *testing.T) {
			importPath := imports.Of(testCase.original)
			assert.Equal(t, testCase.expected, importPath.String())
		})
	}
}

func TestPath_IsEmpty(t *testing.T) {
	testCases := []struct {
		original string
		expected bool
	}{
		{"some/path/to/test/", false},
		{"", true},
		{"/", true},
		{"./", true},
		{"\\", true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.original, func(t *testing.T) {
			importPath := imports.Of(testCase.original)
			assert.Equal(t, testCase.expected, importPath.IsEmpty())
		})
	}
}

func TestPath_Intersection(t *testing.T) {
	testCases := []struct {
		path1    string
		path2    string
		expected string
	}{
		{
			path1:    "/first1/second1/third1/fourth1/fifth1/",
			path2:    "/first/second/third/fourth/fifth/",
			expected: imports.EmptyPath.String(),
		},
		{
			path1:    "/first/second/third/fourth/fifth/",
			path2:    "/first/second/third/fourth/fifth/",
			expected: "first/second/third/fourth/fifth",
		},
		{
			path1:    "/first/second/third",
			path2:    "second/third/fourth/fifth",
			expected: "second/third",
		},
		{
			path1:    "/first/second/third/fourth/fifth/",
			path2:    "/first/second1/third/fourth/fifth/",
			expected: "third/fourth/fifth",
		},
		{
			path1:    "/first/second/third/fourth/fifth/",
			path2:    "/first/second/third1/fourth/fifth/",
			expected: "fourth/fifth",
		},
		{
			path1:    "/first/second/third/fourth/fifth/",
			path2:    "second/third1/fourth/fifth/",
			expected: "fourth/fifth",
		},
		{
			path1:    "/first/second/third/fourth/fifth/",
			path2:    "/first1/second1/third/fourth/fifth/",
			expected: "third/fourth/fifth",
		},
		{
			path1:    "C:\\first\\second\\third\\fourth\\fifth",
			path2:    "C:/first/second/third/fourth/fifth/",
			expected: "C:/first/second/third/fourth/fifth",
		},
		{
			path1:    "C:\\first\\second\\third",
			path2:    "second/third/fourth/fifth",
			expected: "second/third",
		},
		{
			path1:    "C:\\first\\second\\third\\fourth\\fifth/",
			path2:    "C:\\first\\second1\\third\\fourth\\fifth",
			expected: "third/fourth/fifth",
		},
		{
			path1:    "C:\\first\\second\\third\\fourth\\fifth/",
			path2:    "C:\\first\\second\\third1\\fourth\\fifth",
			expected: "C:/first/second",
		},
		{
			path1:    "C:\\first\\second\\third\\fourth\\fifth/",
			path2:    "second/third1/fourth/fifth/",
			expected: "fourth/fifth",
		},
		{
			path1:    "C:\\first\\second\\third\\fourth\\fifth/",
			path2:    "C:\\first1\\second1\\third\\fourth\\fifth",
			expected: "third/fourth/fifth",
		},
		{ // Real case when we have no mod file for sub folder where we run generation
			path1:    "C:\\Users\\username\\goworkspace\\go-annotation\\internal\\lookup\\fixtures",
			path2:    "github.com\\YReshetko\\go-annotation\\internal\\lookup\\fixtures\\dashed-package",
			expected: "go-annotation/internal/lookup/fixtures",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.expected, func(t *testing.T) {
			first := imports.Of(testCase.path1)
			second := imports.Of(testCase.path2)
			assert.Equal(t, filepath.Clean(testCase.expected), first.Intersection(second).String())
			assert.Equal(t, filepath.Clean(testCase.expected), second.Intersection(first).String())
		})
	}
}

func TestPath_Left(t *testing.T) {
	testCases := []struct {
		path1    string
		path2    string
		expected string
	}{
		{
			path1:    "/first/second/third/fourth/fifth/",
			path2:    "/first/second/third/fourth/fifth/",
			expected: imports.EmptyPath.String(),
		},
		{
			path1:    "/first/second/third",
			path2:    "second/third/fourth/fifth",
			expected: "first",
		},
		{
			path1:    "/first/second/third/fourth/fifth/",
			path2:    "/first/second1/third/fourth/fifth/",
			expected: "first/second",
		},
		{
			path1:    "/first/second/third/fourth/fifth/",
			path2:    "/first/second/third1/fourth/fifth/",
			expected: "first/second/third",
		},
		{
			path1:    "/first/second/third/fourth/fifth/",
			path2:    "second/third1/fourth/fifth/",
			expected: "first/second/third",
		},
		{
			path1:    "C:\\first\\second\\third\\fourth\\fifth",
			path2:    "C:/first/second/third/fourth/fifth/",
			expected: imports.EmptyPath.String(),
		},
		{
			path1:    "C:\\first\\second\\third",
			path2:    "second/third/fourth/fifth",
			expected: "C:/first",
		},
		{
			path1:    "C:\\first\\second\\third\\fourth\\fifth/",
			path2:    "C:\\first\\second1\\third\\fourth\\fifth",
			expected: "C:/first/second",
		},
		{
			path1:    "C:\\first\\second\\third\\fourth\\fifth/",
			path2:    "C:\\first\\second\\third1\\fourth\\fifth",
			expected: imports.EmptyPath.String(),
		},
		{
			path1:    "C:\\first\\second\\third\\fourth\\fifth/",
			path2:    "second/third1/fourth/fifth/",
			expected: "C:/first/second/third",
		},
		{ // Real case when we have no mod file for sub folder where we run generation
			path1:    "C:\\Users\\username\\goworkspace\\go-annotation\\internal\\lookup\\fixtures",
			path2:    "github.com\\YReshetko\\go-annotation\\internal\\lookup\\fixtures\\dashed-package",
			expected: "C:/Users/username/goworkspace",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.expected, func(t *testing.T) {
			first := imports.Of(testCase.path1)
			second := imports.Of(testCase.path2)
			assert.Equal(t, filepath.Clean(testCase.expected), first.Left(second).String())
		})
	}
}

func TestPath_Joins(t *testing.T) {
	testCases := []struct {
		path1             string
		path2             string
		expectedLeft      string
		expectedLeftJoin  string
		expectedRight     string
		expectedRightJoin string
		expectedFullJoin  string
	}{
		{
			path1:             "/first/second/third/fourth/fifth/",
			path2:             "/first/second/third/fourth/fifth/",
			expectedLeft:      imports.EmptyPath.String(),
			expectedLeftJoin:  "first/second/third/fourth/fifth",
			expectedRight:     imports.EmptyPath.String(),
			expectedRightJoin: "first/second/third/fourth/fifth",
			expectedFullJoin:  "first/second/third/fourth/fifth",
		},
		{
			path1:             "/first/second/third",
			path2:             "second/third/fourth/fifth",
			expectedLeft:      "first",
			expectedLeftJoin:  "first/second/third",
			expectedRight:     imports.EmptyPath.String(),
			expectedRightJoin: "second/third",
			expectedFullJoin:  "first/second/third/fourth/fifth",
		},
		{
			path1:             "/first/second/third/fourth/fifth/",
			path2:             "/first/second1/third/fourth/fifth/",
			expectedLeft:      "first/second",
			expectedLeftJoin:  "first/second/third/fourth/fifth",
			expectedRight:     imports.EmptyPath.String(),
			expectedRightJoin: "third/fourth/fifth",
			expectedFullJoin:  "first/second/third/fourth/fifth",
		},
		{
			path1:             "/first/second/third/fourth/fifth/",
			path2:             "/first/second/third1/fourth/fifth/",
			expectedLeft:      "first/second/third",
			expectedLeftJoin:  "first/second/third/fourth/fifth",
			expectedRight:     imports.EmptyPath.String(),
			expectedRightJoin: "fourth/fifth",
			expectedFullJoin:  "first/second/third/fourth/fifth",
		},
		{
			path1:             "/first/second/third/fourth/fifth/sixth",
			path2:             "second/third1/fourth/fifth/",
			expectedLeft:      "first/second/third",
			expectedLeftJoin:  "first/second/third/fourth/fifth",
			expectedRight:     "sixth",
			expectedRightJoin: "fourth/fifth/sixth",
			expectedFullJoin:  "first/second/third/fourth/fifth",
		},
		{
			path1:             "second/third1/fourth/fifth/",
			path2:             "/first/second/third/fourth/fifth/sixth",
			expectedLeft:      "second/third1",
			expectedLeftJoin:  "second/third1/fourth/fifth",
			expectedRight:     imports.EmptyPath.String(),
			expectedRightJoin: "fourth/fifth",
			expectedFullJoin:  "second/third1/fourth/fifth/sixth",
		},
		{
			path1:             "C:\\first\\second\\third\\fourth\\fifth",
			path2:             "C:/first/second/third/fourth/fifth/",
			expectedLeft:      imports.EmptyPath.String(),
			expectedLeftJoin:  "C:/first/second/third/fourth/fifth",
			expectedRight:     imports.EmptyPath.String(),
			expectedRightJoin: "C:/first/second/third/fourth/fifth",
			expectedFullJoin:  "C:/first/second/third/fourth/fifth",
		},
		{
			path1:             "C:\\first\\second\\third\\fourth1\\fifth",
			path2:             "second/third/fourth/fifth",
			expectedLeft:      "C:/first",
			expectedLeftJoin:  "C:/first/second/third",
			expectedRight:     "fourth1/fifth",
			expectedRightJoin: "second/third/fourth1/fifth",
			expectedFullJoin:  "C:/first/second/third/fourth/fifth",
		},
		{
			path1:             "C:\\first\\second\\third\\fourth\\fifth/",
			path2:             "C:\\first\\second1\\third\\fourth\\fifth",
			expectedLeft:      "C:/first/second",
			expectedLeftJoin:  "C:/first/second/third/fourth/fifth",
			expectedRight:     imports.EmptyPath.String(),
			expectedRightJoin: "third/fourth/fifth",
			expectedFullJoin:  "C:/first/second/third/fourth/fifth",
		},
		{
			path1:             "C:\\first\\second\\third\\fourth\\fifth/",
			path2:             "C:\\first\\second\\third1\\fourth\\fifth",
			expectedLeft:      imports.EmptyPath.String(),
			expectedLeftJoin:  "C:/first/second",
			expectedRight:     "third/fourth/fifth",
			expectedRightJoin: "C:/first/second/third/fourth/fifth",
			expectedFullJoin:  "C:/first/second/third1/fourth/fifth",
		},
		{
			path1:             "C:\\first\\second\\third\\fourth\\fifth/",
			path2:             "second/third1/fourth/fifth/",
			expectedLeft:      "C:/first/second/third",
			expectedLeftJoin:  "C:/first/second/third/fourth/fifth",
			expectedRight:     imports.EmptyPath.String(),
			expectedRightJoin: "fourth/fifth",
			expectedFullJoin:  "C:/first/second/third/fourth/fifth",
		},
		{ // Real case when we have no mod file for sub folder where we run generation
			path1:             "C:\\Users\\username\\goworkspace\\go-annotation\\internal\\lookup\\fixtures",
			path2:             "github.com\\YReshetko\\go-annotation\\internal\\lookup\\fixtures\\dashed-package",
			expectedLeft:      "C:/Users/username/goworkspace",
			expectedLeftJoin:  "C:/Users/username/goworkspace/go-annotation/internal/lookup/fixtures",
			expectedRight:     imports.EmptyPath.String(),
			expectedRightJoin: "go-annotation/internal/lookup/fixtures",
			expectedFullJoin:  "C:/Users/username/goworkspace/go-annotation/internal/lookup/fixtures/dashed-package",
		},
		{ // Real case when we have no mod file for sub folder where we run generation
			path1:             "github.com\\YReshetko\\go-annotation\\internal\\lookup\\fixtures\\dashed-package",
			path2:             "C:\\Users\\username\\goworkspace\\go-annotation\\internal\\lookup\\fixtures",
			expectedLeft:      "github.com/YReshetko",
			expectedLeftJoin:  "github.com/YReshetko/go-annotation/internal/lookup/fixtures",
			expectedRight:     "dashed-package",
			expectedRightJoin: "go-annotation/internal/lookup/fixtures/dashed-package",
			expectedFullJoin:  "github.com/YReshetko/go-annotation/internal/lookup/fixtures",
		},
	}
	for i, testCase := range testCases {
		first := imports.Of(testCase.path1)
		second := imports.Of(testCase.path2)

		for scenario, fn := range map[[2]string]func(path imports.Path) imports.Path{
			[2]string{"Left", testCase.expectedLeft}:           first.Left,
			[2]string{"Right", testCase.expectedRight}:         first.Right,
			[2]string{"LeftJoin", testCase.expectedLeftJoin}:   first.LeftJoin,
			[2]string{"RightJoin", testCase.expectedRightJoin}: first.RightJoin,
			[2]string{"FullJoin", testCase.expectedFullJoin}:   first.FullJoin,
		} {
			testName := fmt.Sprintf("(%d)_%s_[%s]->[%s]", i, scenario[0], testCase.path1, testCase.path2)
			expected := filepath.Clean(scenario[1])
			t.Run(testName, func(t *testing.T) {
				assert.Equal(t, expected, fn(second).String())
			})
		}
	}
}
