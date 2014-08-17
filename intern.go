// Intern (cache) strings to uint32 values.
package intern

import (
	"sync/atomic"
	"unsafe"
)

// An interned string value.
type S uint32

var internST = unsafe.Pointer(&s0)
var s0 = state{m: map[string]S{"": 0}, r: []string{""}}

type state struct {
	m map[string]S
	r []string
	c S
}

// Intern a string.
func Intern(s string) S {
bset:
	ptr := atomic.LoadPointer(&internST)
	st := (*state)(ptr)
	v, ok := st.m[s]
	if !ok {
		x := *st
		x.m = make(map[string]S, len(st.m))
		for k, v := range st.m {
			x.m[k] = v
		}
		x.c++
		x.m[s] = x.c
		x.r = append(x.r, s)
		if !atomic.CompareAndSwapPointer(&internST, ptr, unsafe.Pointer(&x)) {
			goto bset
		}
	}
	return v
}

// Intern all strings aggregating the write.
func InternAll(ss []string) []S {
	res := make([]S, len(ss))
bset:
	ptr := atomic.LoadPointer(&internST)
	st := (*state)(ptr)
	var x *state
	for i, s := range ss {
		v, ok := st.m[s]
		if ok {
			res[i] = v
			continue
		}
		if x == nil {
			x = new(state)
			*x = *st
			x.m = make(map[string]S, len(st.m))
			for k, v := range st.m {
				x.m[k] = v
			}
		}
		x.c++
		x.m[s] = x.c
		x.r = append(x.r, s)
	}
	if x != nil && !atomic.CompareAndSwapPointer(&internST, ptr, unsafe.Pointer(&x)) {
		goto bset
	}
	return res
}

// Return the string corresponding to an interned string.
func (v S) String() string {
	st := (*state)(atomic.LoadPointer(&internST))
	return st.r[v]
}
