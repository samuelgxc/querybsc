package exchange

import "querybsc/common/db"

var exchange map[string]string

func Start() {
	db.Session().Select("address", "name").From("exchange").Load(&exchange)

}
func Get(address string, lable string) string {
	if newlable, ok := exchange[address]; ok {
		return newlable
	} else {
		return lable
	}
}

func Have(address string) bool {
	_, ok := exchange[address]
	return ok
}
