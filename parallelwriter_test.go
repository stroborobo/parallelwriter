package parallelwriter

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestParallelWriter(t *testing.T) {
	b1 := &bytes.Buffer{}
	b2 := &bytes.Buffer{}

	pw := ParallelWriter(b1, b2)

	if _, ok := pw.(*parallelWriter); !ok {
		t.Error("ParallelWriter return is not a parallelWriter!")
	}

}

func TestWrite(t *testing.T) {
	b1 := &bytes.Buffer{}
	b2 := &bytes.Buffer{}

	pw := ParallelWriter(b1, b2)

	data := []byte("foobar")
	n, err := pw.Write(data)
	if err != nil {
		t.Error("Error on write:", err)
	}
	if n != len(data) {
		t.Error("Written count is not equal to data length:", len(data), "!=", n)
	}

	if b := b1.Bytes(); !bytes.Equal(b, data) {
		t.Error("b1 content is incorrect. Expected:", data, "Got:", b)
	}
	if b := b2.Bytes(); !bytes.Equal(b, data) {
		t.Error("b2 content is incorrect. Expected:", data, "Got:", b)
	}
}

type slowBuffer struct {
	buf bytes.Buffer
}

func (s *slowBuffer) Write(data []byte) (int, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	<-time.After(time.Millisecond * time.Duration(rand.Int()%1000))
	return s.buf.Write(data)
}

func (s *slowBuffer) Bytes() []byte {
	return s.buf.Bytes()
}

func TestSlowWrite(t *testing.T) {
	b1 := &slowBuffer{}
	b2 := &slowBuffer{}

	pw := ParallelWriter(b1, b2)

	data := []byte("foobar")
	n, err := pw.Write(data)
	if err != nil {
		t.Error("Error on write:", err)
	}
	if n != len(data) {
		t.Error("Written count is not equal to data length:", len(data), "!=", n)
	}

	if b := b1.Bytes(); !bytes.Equal(b, data) {
		t.Error("b1 content is incorrect. Expected:", data, "Got:", b)
	}
	if b := b2.Bytes(); !bytes.Equal(b, data) {
		t.Error("b2 content is incorrect. Expected:", data, "Got:", b)
	}
}
