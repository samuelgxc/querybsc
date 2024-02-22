package page

import (
	"fmt"
	"shangchenggo/common/db"
	"testing"
)

func TestNewPaginator(t *testing.T) {
	req := TestReq{}
	req.Page = 1
	req.Limit = 3
	res := ServiceFunc(req)
	fmt.Println(res)
}

func ServiceFunc(req TestReq) TestRes {
	res := TestRes{}
	stmt := db.Session().Select("id, username, nickname, guid, guid_time, approve_user").From("user").OrderAsc("id")
	pdata, err := NewPaginator(stmt, req.Page, req.Limit)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(pdata)
		err = pdata.Load(&res)
		if err != nil {
			fmt.Println(err)
		}
	}
	return res
}

type TestReq struct {
	PageReq
}

type TestRes struct {
	PageRes
	Data []DataRes
}

type DataRes struct {
	Id          int64
	Username    string
	Nickname    string
	Guid        string
	GuidTime    string
	ApproveUser int64
}
