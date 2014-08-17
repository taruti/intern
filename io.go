package intern

import (
	"bufio"
	"io"
	"sync/atomic"
	"unsafe"
)

type sbwriter interface {
	WriteByte(byte) error
	WriteString(string) (int, error)
}

func (ctx *Context) WriteTo(rawwr io.Writer) error {
	w, ok := rawwr.(sbwriter)
	var bwr *bufio.Writer
	if !ok {
		bwr = bufio.NewWriter(rawwr)
		w = bwr
	}
	c := (*state)(atomic.LoadPointer(&ctx.p))
	for _, s := range c.r {
		_, e := w.WriteString(s)
		if e != nil {
			return e
		}
		e = w.WriteByte('\n')
		if e != nil {
			return e
		}
	}
	if bwr != nil {
		return bwr.Flush()
	}
	return nil
}

func ReadContext(rawrd io.Reader) (Context, error) {
	rd := bufio.NewReader(rawrd)
	st := newst()
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				return Context{unsafe.Pointer(st)}, nil
			}
			return Context{}, err
		}
		st.addMissing(string(line))
	}
}
