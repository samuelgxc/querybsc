package valid

type Func struct {
	name string
	msg  string
	fn   func(CV) bool
}

func (r *Func) Valid(m map[string]string) (string, bool) {
	v, _ := m[r.name]
	if !r.fn(CV{Value: v, Map: m}) {
		return r.msg, true
	} else {
		return "", false
	}
}
func (r *Func) GetName() string {
	return r.name
}
