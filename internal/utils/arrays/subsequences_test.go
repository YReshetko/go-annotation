package arrays_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/YReshetko/go-annotation/internal/utils/arrays"
)

func TestLCS_Integers(t *testing.T) {
	a := []int{1, 2, 8, 2, 1, 4, 6, 1, 4, 8, 2}
	b := []int{6, 8, 2, 1, 5, 4, 1, 8, 2}

	expected := []int{8, 2, 1}
	expectedI := 2
	expectedJ := 1

	actual, actualI, actualJ := arrays.LCS(a, b)
	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedI, actualI)
	assert.Equal(t, expectedJ, actualJ)
}

func TestLCS_Strings(t *testing.T) {
	a := []string{"C:", "Users", "username", "goworkspace", "go-annotation", "internal", "lookup", "fixtures"}
	b := []string{"github.com", "YReshetko", "go-annotation", "internal", "lookup", "fixtures", "dashed-package"}

	expected := []string{"go-annotation", "internal", "lookup", "fixtures"}
	expectedI := 4
	expectedJ := 2

	actual, actualI, actualJ := arrays.LCS(a, b)
	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedI, actualI)
	assert.Equal(t, expectedJ, actualJ)
}

func TestLCS_NoSubsequence(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}

	var expected []int
	expectedI := -1
	expectedJ := -1

	actual, actualI, actualJ := arrays.LCS(a, b)
	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedI, actualI)
	assert.Equal(t, expectedJ, actualJ)
}
