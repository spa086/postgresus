package usecases_postgresql

import "io"

// CountingWriter wraps an io.Writer and counts the bytes written to it
type CountingWriter struct {
	writer       io.Writer
	bytesWritten int64
}

func (cw *CountingWriter) Write(p []byte) (n int, err error) {
	n, err = cw.writer.Write(p)
	cw.bytesWritten += int64(n)
	return n, err
}

// GetBytesWritten returns the total number of bytes written
func (cw *CountingWriter) GetBytesWritten() int64 {
	return cw.bytesWritten
}
