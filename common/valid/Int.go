package valid

import val "github.com/asaskevich/govalidator"

type Int struct {
	name string
	msg  string
}

func (r *Int) Valid(m map[string]string) (string, bool) {
	v, _ := m[r.name]
	if !val.IsInt(v) {
		//return
		return r.msg, true
	} else {
		return "", false
	}
}
func (r *Int) GetName() string {
	return r.name
}
