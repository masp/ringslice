package ring

import (
	"reflect"
	"slices"
	"testing"
)

func array[T any](r Slice[T]) []T {
	return slices.Collect(r.Values())
}

func TestAppend(t *testing.T) {
	r := NewSlice([]int{1, 2, 3, 4, 5})
	r.Append(6, 7, 8)
	if !reflect.DeepEqual(array(r), []int{6, 7, 8}) {
		t.Errorf("expected [6 7 8], got %v", array(r))
	}
	if !reflect.DeepEqual(r.Slice(), []int{6, 7, 8}) {
		t.Errorf("expected [6 7], got %v", r.Slice())
	}

	r.Append(9, 10, 11, 12, 13, 14, 15)
	if !reflect.DeepEqual(array(r), []int{11, 12, 13, 14, 15}) {
		t.Errorf("expected [11 12 13 14 15], got %v", array(r))
	}
	if !reflect.DeepEqual(r.Slice(), []int{11, 12, 13, 14, 15}) {
		t.Errorf("expected [11 12 13 14 15], got %v", array(r))
	}
}

func TestCopy(t *testing.T) {
	r := NewSlice(make([]int, 3))
	r.Append(6, 7)
	dst := make([]int, 5)
	n := r.CopyTo(dst)
	if !reflect.DeepEqual(dst, []int{6, 7, 0, 0, 0}) {
		t.Errorf("expected [6 7 0 0 0], got %v", dst)
	}
	if n != 2 {
		t.Errorf("expected n=2, got %v", n)
	}

	dst = make([]int, 2)
	n = r.CopyTo(dst)
	if !reflect.DeepEqual(dst, []int{6, 7}) {
		t.Errorf("expected [6 7], got %v", dst)
	}
	if n != 2 {
		t.Errorf("expected n=2, got %v", n)
	}

	dst = make([]int, 1)
	n = r.CopyTo(dst)
	if !reflect.DeepEqual(dst, []int{6}) {
		t.Errorf("expected [6], got %v", dst)
	}
	if n != 1 {
		t.Errorf("expected n=1, got %v", n)
	}
}

func TestBuffer(t *testing.T) {
	buf := NewBuffer(make([]byte, 3))
	buf.Write([]byte("abc123xy"))
	if buf.String() != "3xy" {
		t.Errorf("expected 3xy, got %v", buf.String())
	}
}
