package algodb

import (
	"errors"
	"io"
)

var (
	errInvalidArgument  = errors.New("invalid argument")
	errNegativePosition = errors.New("negative position")
	errNegativeOffset   = errors.New("negative offset")
)

func checkNegativeOffset(off int64) (err error) {
	if off < 0 {
		err = errNegativeOffset
	}
	return
}

type ByteBuffer struct {
	off int64
	b   []byte
}

func NewByteBuffer(b []byte) *ByteBuffer {
	return &ByteBuffer{
		b: b,
	}
}

func (b *ByteBuffer) writeAt(p []byte, off int64) (n int, err error) {
	prevLen := len(b.b)
	diff := int(off) - prevLen
	if diff > 0 {
		b.b = append(b.b, make([]byte, diff)...)
	}
	b.b = append(b.b[:off], p...)
	if len(b.b) < prevLen {
		b.b = b.b[:prevLen]
	}
	return len(p), nil
}

func (b *ByteBuffer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	n, err = b.writeAt(p, b.off)
	if err != nil {
		return
	}
	b.off += int64(n)
	return
}

func (b *ByteBuffer) WriteAt(p []byte, off int64) (n int, err error) {
	if err = checkNegativeOffset(off); err != nil {
		return
	}
	if len(p) == 0 {
		return
	}
	n, err = b.writeAt(p, off)
	return
}

func (b *ByteBuffer) Seek(offset int64, whence int) (pos int64, err error) {
	switch whence {
	case io.SeekStart:
		pos = offset
	case io.SeekCurrent:
		pos = offset + b.off
	case io.SeekEnd:
		pos = offset + int64(len(b.b))
	default:
		err = errInvalidArgument
		return
	}
	if pos < 0 {
		err = errNegativePosition
		return
	}
	b.off = pos
	return
}

func (b *ByteBuffer) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	if b.off >= int64(len(b.b)) {
		return 0, io.EOF
	}
	n = copy(p, b.b[b.off:])
	b.off += int64(n)
	return
}

func (b *ByteBuffer) ReadAt(p []byte, off int64) (n int, err error) {
	if err = checkNegativeOffset(off); err != nil {
		return
	}
	if len(p) == 0 {
		return
	}
	if off >= int64(len(b.b)) {
		return 0, io.EOF
	}
	n = copy(p, b.b[off:])
	if n < len(p) {
		err = io.EOF
	}
	return
}

func (b *ByteBuffer) Len() int {
	return len(b.b)
}

func (b *ByteBuffer) Bytes() []byte {
	return b.b
}
