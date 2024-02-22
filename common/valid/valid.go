package valid

import (
	val "github.com/asaskevich/govalidator"
	"strconv"
	//"fmt"
)

type Vdate int64

/*
self::EXISTS_VALIDATE 或者0 存在字段就验证（默认）
self::MUST_VALIDATE 或者1 必须验证
self::VALUE_VALIDATE或者2 值不为空的时候验证
*/
const EXISTS = Vdate(0) //存在字段就验证
const MUST = Vdate(1)   //必须验证
const VALUE = Vdate(2)  //值不为空的时候验证

type ValidInterface interface {
	Valid(map[string]string) (string, bool)
	GetName() string
}

type Valid struct {
	rule  []ValidInterface
	calls []IfCall
}

//check val
type CV struct {
	Value string
	Map   map[string]string
}

func (c *CV) Get(n string) string {
	ret, _ := c.Map[n]
	return ret
}
func (v *Valid) Check(m map[string]string) (ret map[string]string) {
	ret = make(map[string]string)
	//对常规验证进行处理
	for _, r := range v.rule {
		msg, find := r.Valid(m)
		if find {
			if _, find := ret[r.GetName()]; !find {
				ret[r.GetName()] = msg
			}
		}
	}
	//对ifcall做处理
	for _, c := range v.calls {
		//触发判定条件
		if c.f() {
			for ckey, cmsg := range c.valid.Check(m) {
				if _, find := ret[ckey]; !find {
					ret[ckey] = cmsg
				}
			}
		}
	}
	return
}

//必填
func (v *Valid) Require(name string, msg string) *Valid {
	v.rule = append(v.rule, &Require{name: name, msg: msg})
	return v
}

func (v *Valid) IfCall(f func() bool, valid Valid) *Valid {
	v.calls = append(v.calls, IfCall{f: f, valid: valid})
	return v
}
func (v *Valid) Func(f string, fn func(CV) bool, msg string) *Valid {
	v.rule = append(v.rule, &Func{name: f, fn: fn, msg: msg})
	return v
}
func (v *Valid) Length(f string, min int, max int, msg string) *Valid {
	v.rule = append(v.rule, &Length{name: f, min: min, max: max, msg: msg})
	return v
}
func (v *Valid) Confirm(f string, c string, msg string) *Valid {
	v.rule = append(v.rule, &Confirm{name: f, cname: c, msg: msg})
	return v
}
func (v *Valid) Regex(f string, exp string, msg string) *Valid {
	v.rule = append(v.rule, &Regex{name: f, exp: exp, msg: msg})
	return v
}
func (v *Valid) Int(f string, msg string) *Valid {
	v.rule = append(v.rule, &Int{name: f, msg: msg})
	return v
}
func (v *Valid) InSlice(f string, s []string, msg string) *Valid {
	v.rule = append(v.rule, &InSlice{name: f, slice: s, msg: msg})
	return v
}
func (v *Valid) Number(f string, msg string) *Valid {
	v.rule = append(v.rule, &Number{name: f, msg: msg})
	return v
}

func (v *Valid) Mx(f string, limit string, limitnum float64, msg string) *Valid {
	v.rule = append(v.rule, &Mx{name: f, msg: msg, limitNum: limitnum, limit: limit})
	return v
}

//正整数验证
func ValidateInt(s string) bool {
	if !val.IsInt(s) {
		return false
	}
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return false
	}
	if i <= 0 {
		return false
	}
	return true
}
