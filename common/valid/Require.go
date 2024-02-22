package valid

type Require struct {
	name string
	msg  string
}

func (r *Require) Valid(m map[string]string) (string, bool) {
	v, _ := m[r.name]
	if v == "" {
		return r.msg, true
	} else {
		return "", false
	}
}
func (r *Require) GetName() string {
	return r.name
}
