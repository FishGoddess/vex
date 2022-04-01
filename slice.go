package vex

// makeBytes makes a new byte slice.
func makeBytes(cap int32) []byte {
	return make([]byte, cap)
}
