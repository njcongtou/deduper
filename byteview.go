package deduper

// ByteView holds an immutable view of bytes.
type ByteView struct {
	b []byte // able to support any type in Go
}

func (bv ByteView) Len() int {
	return len(bv.b)
}

// ByteSlice returns a copy of the data a byte slice.
// This avoids the actual can be modified from client
func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

// String returns the data as a string, making a copy if necessary.
func (bv ByteView) String() string {
	return string(bv.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
