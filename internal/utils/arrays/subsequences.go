package arrays

// LCS - Longest Common Subsequence
// The function searches for LCS and returns it along with first index in the slice 'a' and the first in the slice 'b'
func LCS[T comparable](a, b []T) ([]T, int, int) {
	dp := make([][]int, len(a))
	var maxValue, maxI, maxJ int

	for i := len(a) - 1; i >= 0; i-- {
		dp[i] = make([]int, len(b))
		for j := len(b) - 1; j >= 0; j-- {
			if a[i] != b[j] {
				continue
			}

			value := 1
			if i != len(a)-1 && j != len(b)-1 {
				value = dp[i+1][j+1] + 1
			}

			dp[i][j] = value

			if value > maxValue {
				maxValue = value
				maxI = i
				maxJ = j
			}
		}
	}

	if maxValue == 0 {
		return nil, -1, -1
	}
	return a[maxI : maxI+maxValue], maxI, maxJ
}
