package algodb

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestNewByteBuffer(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want *ByteBuffer
	}{
		{
			name: "nil",
			want: &ByteBuffer{},
		},
		{
			name: "empty",
			args: args{
				b: []byte{},
			},
			want: &ByteBuffer{
				b: []byte{},
			},
		},
		{
			name: "0123456789",
			args: args{
				b: []byte("0123456789"),
			},
			want: &ByteBuffer{
				b: []byte("0123456789"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewByteBuffer(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewByteBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByteBufferWrite(t *testing.T) {
	var buf ByteBuffer
	tests := []struct {
		p []byte
	}{
		{p: nil},
		{p: []byte{1}},
		{p: []byte{0, 1, 0}},
	}
	for _, tt := range tests {
		n, err := buf.Write(tt.p)
		if err != nil {
			t.Errorf("error = %v, want = nil", err)
			return
		}
		if n != len(tt.p) {
			t.Errorf("n = %d, want %d", n, len(tt.p))
		}
	}
	wantBytes := []byte{1, 0, 1, 0}
	if gotBytes := buf.Bytes(); !bytes.Equal(wantBytes, gotBytes) {
		t.Errorf("ByteBuffer.Bytes() = %v, want %v", gotBytes, wantBytes)
	}
}

func TestByteBufferWriteAtNegativeOffset(t *testing.T) {
	var buf ByteBuffer
	n, err := buf.WriteAt([]byte{}, -1)
	if err == nil {
		t.Errorf("error = %v, want error negative offset", n)
		return
	}
	if n != 0 {
		t.Errorf("n = %d, want 0", n)
	}
}

func TestByteBufferWriteAt(t *testing.T) {
	type args struct {
		p   []byte
		off int64
	}
	var buf ByteBuffer
	tests := []struct {
		args      args
		wantN     int
		wantBytes []byte
	}{
		{
			args: args{
				p:   nil,
				off: 1,
			},
			wantN:     0,
			wantBytes: nil,
		},
		{
			args: args{
				p:   []byte{},
				off: 1,
			},
			wantN:     0,
			wantBytes: nil,
		},
		{
			args: args{
				p:   []byte("89"),
				off: 8,
			},
			wantN:     2,
			wantBytes: []byte{0, 0, 0, 0, 0, 0, 0, 0, 56, 57},
		},
		{
			args: args{
				p:   []byte("0123"),
				off: 0,
			},
			wantN:     4,
			wantBytes: []byte{48, 49, 50, 51, 0, 0, 0, 0, 56, 57},
		},
		{
			args: args{
				p:   []byte("4567"),
				off: 4,
			},
			wantN:     4,
			wantBytes: []byte("0123456789"),
		},
	}
	for _, tt := range tests {
		n, err := buf.WriteAt(tt.args.p, tt.args.off)
		if err != nil {
			t.Errorf("error = %v, want = nil", err)
		}
		if n != tt.wantN {
			t.Errorf("n = %d, want %d", n, tt.wantN)
		}
		if gotBytes := buf.Bytes(); !reflect.DeepEqual(gotBytes, tt.wantBytes) {
			t.Errorf("bytes = %v, want %v", gotBytes, tt.wantBytes)
		}
	}
}

func TestByteBufferSeek(t *testing.T) {
	type args struct {
		b      []byte
		offset int64
		whence int
	}
	tests := []struct {
		name    string
		args    args
		wantPos int64
		wantErr bool
	}{
		{
			name: "InvalidWhence",
			args: args{
				whence: 3,
			},
			wantErr: true,
			wantPos: 0,
		},
		{
			name: "NegativePosition",
			args: args{
				whence: io.SeekStart,
				offset: -10,
			},
			wantPos: -10,
			wantErr: true,
		},
		{
			name: "Start",
			args: args{
				b:      []byte("0123456789"),
				offset: 0,
				whence: io.SeekStart,
			},
			wantPos: 0,
		},
		{
			name: "Current",
			args: args{
				b:      []byte("0123456789"),
				offset: 1,
				whence: io.SeekCurrent,
			},
			wantPos: 1,
		},
		{
			name: "End",
			args: args{
				b:      []byte("0123456789"),
				offset: 0,
				whence: io.SeekEnd,
			},
			wantPos: 10,
		}, {
			name: "NegativeOffset",
			args: args{
				b:      []byte("0123456789"),
				offset: -1,
				whence: io.SeekEnd,
			},
			wantPos: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := NewByteBuffer(tt.args.b)
			gotPos, err := buf.Seek(tt.args.offset, tt.args.whence)
			if (err != nil) != tt.wantErr {
				t.Errorf("Seek() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPos != tt.wantPos {
				t.Errorf("Seek() gotPos = %v, want %v", gotPos, tt.wantPos)
			}
		})
	}
}

func TestByteBufferRead(t *testing.T) {
	buf := NewByteBuffer([]byte("0123456789"))
	tests := []struct {
		size  int
		wantN int
		want  []byte
	}{
		{
			size:  0,
			wantN: 0,
			want:  []byte{},
		},
		{
			size:  4,
			wantN: 4,
			want:  []byte("0123"),
		},
		{
			size:  4,
			wantN: 4,
			want:  []byte("4567"),
		},
		{
			size:  2,
			wantN: 2,
			want:  []byte("89"),
		},
	}
	for _, tt := range tests {
		p := make([]byte, tt.size)
		n, err := buf.Read(p)
		if err != nil {
			t.Errorf("error = %v, want = nil", err)
			return
		}
		if n != tt.wantN {
			t.Errorf("n = %d, want %d", n, tt.wantN)
			return
		}
		if !reflect.DeepEqual(p, tt.want) {
			t.Errorf("p = %v, want %v", p, tt.want)
		}
	}
}

func TestByteBufferReadAll(t *testing.T) {
	buf := NewByteBuffer([]byte("0123456789"))
	gotBytes, err := io.ReadAll(buf)
	if err != nil {
		t.Errorf("error = %v, want nil", err)
		return
	}
	wantBytes := []byte("0123456789")
	if !bytes.Equal(wantBytes, gotBytes) {
		t.Errorf("io.ReadAll() = %v, want %v", gotBytes, wantBytes)
	}
}

func TestByteBufferReadAtNegativeOffset(t *testing.T) {
	var buf ByteBuffer
	n, err := buf.ReadAt([]byte{}, -1)
	if err == nil {
		t.Errorf("error = %v, want error negative offset", n)
		return
	}
	if n != 0 {
		t.Errorf("n = %d, want 0", n)
	}
}

func TestByteBufferReadAt(t *testing.T) {
	buf := NewByteBuffer([]byte("0123456789"))
	tests := []struct {
		off       int64
		size      int
		wantN     int
		wantErr   bool
		wantBytes []byte
	}{
		{
			size:      0,
			wantN:     0,
			wantBytes: []byte{},
		},
		{
			size:      10,
			wantN:     10,
			wantBytes: []byte("0123456789"),
		},
		{
			size:      4,
			off:       4,
			wantN:     4,
			wantBytes: []byte("4567"),
		},
		{
			size:      4,
			off:       8,
			wantN:     2,
			wantErr:   true,
			wantBytes: []byte{56, 57, 0, 0},
		},
	}
	for _, tt := range tests {
		p := make([]byte, tt.size)
		n, err := buf.ReadAt(p, tt.off)
		if (err != nil) != tt.wantErr {
			t.Errorf("error = %v, want = nil", err)
			return
		}
		if n != tt.wantN {
			t.Errorf("n = %d, want %d", n, tt.wantN)
			return
		}
		if !reflect.DeepEqual(p, tt.wantBytes) {
			t.Errorf("p = %v, want %v", p, tt.wantBytes)
		}
	}
}
