// Intern (cache) strings to uint32 values.
package intern

import (
	"sync/atomic"
	"unsafe"
)

// An interned string value from some Context.
type S uint32

type state struct {
	m map[string]S
	r []string
	c S
}

type Context struct {
	p unsafe.Pointer
}

// Intern a string.
func (c *Context) Intern(s string) S {
bset:
	ptr := atomic.LoadPointer(&c.p)
	st := (*state)(ptr)
	v, ok := st.m[s]
	if !ok {
		x := *st
		x.m = make(map[string]S, len(st.m))
		for k, v := range st.m {
			x.m[k] = v
		}
		x.addMissing(s)
		if !atomic.CompareAndSwapPointer(&c.p, ptr, unsafe.Pointer(&x)) {
			goto bset
		}
	}
	return v
}

// Intern all strings aggregating the write.
func (c *Context) InternAll(ss []string, mayberes []S) {
bset:
	ptr := atomic.LoadPointer(&c.p)
	st := (*state)(ptr)
	var newstptr *state
	for i, s := range ss {
		v, ok := st.m[s]
		if ok {
			if mayberes != nil {
				mayberes[i] = v
			}
			continue
		}
		if newstptr == nil {
			newstptr = new(state)
			*newstptr = *st
			newstptr.m = make(map[string]S, len(st.m))
			for k, v := range st.m {
				newstptr.m[k] = v
			}
		}
		newstptr.addMissing(s)
	}
	if newstptr != nil && !atomic.CompareAndSwapPointer(&c.p, ptr, unsafe.Pointer(newstptr)) {
		goto bset
	}
}

func (st *state) addMissing(s string) {
	st.c++
	st.m[s] = st.c
	st.r = append(st.r, s)
}

// Return the string corresponding to an interned string.
func (c *Context) String(v S) string {
	st := (*state)(atomic.LoadPointer(&c.p))
	return st.r[v]
}

// Create a new Context.
func NewContext() Context {
	return Context{unsafe.Pointer(newst())}
}

func newst() *state {
	const s = 4096
	mm := make(map[string]S, s)
	mm[""] = 0
	return &state{m: mm, r: make([]string, 1, s)}
}
