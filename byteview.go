package naive_cache

// ByteView holds an immutable view of bytes
type  ByteView struct {
	b []byte
}

// Len return the view's len
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice returns a copy slice of byte
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// String return the data as string
func (v ByteView) String() string {
	return string(v.b)
}

