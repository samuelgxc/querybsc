package valid

import (
	val "github.com/asaskevich/govalidator"
)

type Mx struct {
	limit    string //max min
	limitNum float64
	name     string
	msg      string
}

func (r *Mx) Valid(m map[string]string) (string, bool) {
	v, _ := m[r.name]
	if !val.IsNumeric(v) {
		//return
		return r.msg, true
	} else {
		l, _ := val.ToFloat(v)
		switch r.limit {
		case "max":
			if l > r.limitNum {
				return r.msg, true
			}
		case "min":
			if l < r.limitNum {
				return r.msg, true
			}
		default:
			return r.msg, true
		}
		return "", false
	}
}
func (r *Mx) GetName() string {
	return r.name
}
