package utils

// Abs returns the absolute value of an integer.
// If the input is negative, it returns the positive equivalent.
// Otherwise, it returns the input as is.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
