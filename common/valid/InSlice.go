package valid

import (
	"fmt"
	"github.com/astaxie/beego/utils"
)

type InSlice struct {
	name  string
	slice []string
	msg   string
}

func (r *InSlice) Valid(m map[string]string) (string, bool) {
	v, _ := m[r.name]
	//fmt.Println(m)
	if !utils.InSlice(v, r.slice) {
		fmt.Println(r.msg)
		return r.msg, true
	} else {
		return "", false
	}
}
func (r *InSlice) GetName() string {
	return r.name
}
