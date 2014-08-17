package intern

import "unsafe"

// Interned string in the global intern table.
type Global S

var internST = NewContext()

// Intern a string in the global table.
func Intern(s string) Global {
	return Global(internST.Intern(s))
}

// Intern strings in the global table.
func InternAll(ss []string, mayberes []Global) {
	if mayberes == nil {
		internST.InternAll(ss, nil)
	} else {
		res := *(*[]S)(unsafe.Pointer(&mayberes))
		internST.InternAll(ss, res)
	}
}

func (g Global) String() string {
	return internST.String(S(g))
}
