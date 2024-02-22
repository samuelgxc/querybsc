package valid

type IfCall struct {
	//name string
	//msg string
	f     func() bool
	valid Valid
}

//func (r *IfCall)Valid(m map[string]string)(string,bool){
//v,_:=m[r.name];
//if v ==""{
//return r.msg,true
//}else{
//	return "",false
//}
//}
