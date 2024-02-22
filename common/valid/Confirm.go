package valid

import "fmt"

type Confirm struct {
	name  string
	cname string
	msg   string
}

func (r *Confirm) Valid(m map[string]string) (string, bool) {
	v, _ := m[r.name]
	c, _ := m[r.cname]

	fmt.Println(r.name, r.cname)
	fmt.Println(v, c)
	if v != c {
		return r.msg, true
	} else {
		return "", false
	}
}
func (r *Confirm) GetName() string {
	return r.name
}
