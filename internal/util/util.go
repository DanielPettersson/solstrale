package util

import (
	"math"
	"os"
)

var (
	// Infinity is the largest possible float64 value
	Infinity float64 = math.Inf(1)
	// AlmostZero is a value that is so small as to be almost zero
	AlmostZero float64 = 1e-8
)

// DegreesToRadians converts an angle in degrees to radians
func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

// IsExecOwner checks the file mode to see if it is executable by owner
func IsExecOwner(mode os.FileMode) bool {
	return mode&0100 != 0
}
