package valid

import val "github.com/asaskevich/govalidator"

type Number struct {
	name string
	msg  string
}

func (r *Number) Valid(m map[string]string) (string, bool) {
	v, _ := m[r.name]
	if !val.IsNumeric(v) {
		//return
		return r.msg, true
	} else {
		return "", false
	}
}
func (r *Number) GetName() string {
	return r.name
}
