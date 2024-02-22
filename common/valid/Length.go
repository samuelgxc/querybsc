package valid

type Length struct {
	name string
	min  int
	max  int
	msg  string
}

func (r *Length) Valid(m map[string]string) (string, bool) {
	v, _ := m[r.name]
	ru := []rune(v)
	if len(ru) < r.min || len(ru) > r.max {
		//return
		return r.msg, true
	} else {
		return "", false
	}
}
func (r *Length) GetName() string {
	return r.name
}
