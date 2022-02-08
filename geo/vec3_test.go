package geo

import "testing"

func TestZeroVector(t *testing.T) {

	len := ZeroVector.Length()

	if len != 0 {
		t.Errorf("Length of zero vector should be 0. Is %v", len)
	}
}
