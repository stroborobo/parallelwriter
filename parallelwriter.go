package parallelwriter

import "io"

type parallelWriter struct {
	writers []io.Writer
}

func (p *parallelWriter) Write(data []byte) (int, error) {
	type ret struct {
		n   int
		err error
	}
	rets := make(chan *ret)

	for _, w := range p.writers {
		go func(writer io.Writer) {
			n, err := writer.Write(data)
			rets <- &ret{n, err}
		}(w)
	}

	for range p.writers {
		r := <-rets
		if r.err != nil {
			return r.n, r.err
		}
		if r.n != len(data) {
			return r.n, io.ErrShortWrite
		}
	}
	return len(data), nil
}

// ParallelWriter creates a writer like io.MultiWriter but executes the writes
// concurrently.
func ParallelWriter(writers ...io.Writer) io.Writer {
	w := make([]io.Writer, len(writers))
	copy(w, writers)
	return &parallelWriter{w}
}
