package ring

import (
	"iter"
)

type Slice[T any] struct {
	data []T
	pos  uint64
}

// New creates a new ring buffer with the given capacity. If len(buf) == 0, NewSlice panics.
func NewSlice[T any](buf []T) Slice[T] {
	if len(buf) == 0 {
		panic("ring buffer must have a size greater than zero")
	}
	return Slice[T]{
		data: buf,
		pos:  0,
	}
}

// Append adds elements to the buffer. If the buffer would grow past capacity, the oldest elements are dropped.
func (cs *Slice[T]) Append(vs ...T) {
	if len(vs) > len(cs.data) {
		vs = vs[len(vs)-len(cs.data):]
	}
	p := cs.pos % uint64(len(cs.data))
	n := copy(cs.data[p:], vs)
	copy(cs.data, vs[n:])
	cs.pos += uint64(len(vs))
}

// CopyTo copies the elements of the buffer to dst and returns the number of elements copied.
func (cs Slice[T]) CopyTo(dst []T) int {
	p := int(cs.pos % uint64(len(cs.data)))
	if cs.pos < uint64(len(cs.data)) {
		return copy(dst, cs.data[:p])
	}

	var n1, n2 int
	n1 = copy(dst, cs.data[p:])
	if n1 < len(dst) {
		n2 = copy(dst[n1:], cs.data[:p])
	}
	return n1 + n2
}

// Values returns the values appended to this slice in insertion order.
func (cs Slice[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		p := int(cs.pos % uint64(len(cs.data)))

		// If cs.pos < len(cs.data), then the buffer isn't full yet.
		if cs.pos < uint64(len(cs.data)) {
			for i := 0; i < p; i++ {
				if !yield(cs.data[i]) {
					return
				}
			}
			return
		}

		// Oldest (from p to len(cs.data))
		for i := p; i < len(cs.data); i++ {
			if !yield(cs.data[i]) {
				return
			}
		}

		// Newest elements (from 0 to p)
		for i := 0; i < p; i++ {
			if !yield(cs.data[i]) {
				return
			}
		}
	}
}

// Slice returns a newly allocated slice containing the last inserted elements.
func (cs Slice[T]) Slice() []T {
	tmp := make([]T, len(cs.data))
	n := cs.CopyTo(tmp)
	return tmp[:n]
}

type Buffer struct {
	Slice[byte]
}

func NewBuffer(buf []byte) Buffer {
	return Buffer{NewSlice(buf)}
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.Append(p...)
	return len(p), nil
}

func (b Buffer) String() string {
	return string(b.Slice.Slice())
}
