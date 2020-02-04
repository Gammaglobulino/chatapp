package trace

import (
	"fmt"
	"io"
)

type Tracer interface {
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}
type nilTracer struct{}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func (nt *nilTracer) Trace(a ...interface{}) {}
func Off() Tracer {
	return &nilTracer{}
}
