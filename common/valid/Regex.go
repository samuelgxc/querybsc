package valid

import "regexp"

type Regex struct {
	name string
	exp  string
	msg  string
}

func (r *Regex) Valid(m map[string]string) (string, bool) {
	v, _ := m[r.name]
	ret, _ := regexp.MatchString(r.exp, v)
	if !ret {
		return r.msg, true
	} else {
		return "", false
	}
}
func (r *Regex) GetName() string {
	return r.name
}
