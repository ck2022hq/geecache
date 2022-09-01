package geecache

import "io"

// A ByteView holds an immutable view of bytes.
type ByteView struct {
	b []byte
}

func (byteView ByteView) Len() int {
	return len(byteView.b)
}

func (byteView ByteView) ByteSlice() []byte {
	return cloneBytes(byteView.b)
}

func cloneBytes(b []byte) []byte {
	result := make([]byte, len(b))
	copy(result, b)
	return result
}

func (byteView ByteView) String() string {
	return string(byteView.b)
}

// WriteTo implements io.WriterTo on the bytes in v.
func (v ByteView) WriteTo(w io.Writer) (n int64, err error) {
	m, err := w.Write(v.b)
	if err == nil && m < v.Len() {
		err = io.ErrShortWrite
	}
	n = int64(m)
	return
}
