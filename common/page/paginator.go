package page

import (
	"errors"
	"github.com/gocraft/dbr"
	"shangchenggo/tools/scan"
	"reflect"
)

//分页器
type Paginator struct {
	Page  uint64 //当前页码
	Limit uint64 //每页条数
	Total uint64 //总条数
	Data  []map[string]interface{}
}

//创建分页器
func NewPaginator(stmt *dbr.SelectStmt, page uint64, limit uint64) (*Paginator, error) {
	p := new(Paginator)
	if page == 0 {
		page = 1
	}
	p.Page = page
	if limit == 0 {
		limit = 10
	}
	p.Limit = limit
	//查询总条数
	rows, err := stmt.Rows()
	if err != nil {
		return nil, err
	}
	var total uint64
	for rows.Next() {
		total++
	}
	p.Total = total
	//查询分页数据
	p.Data = []map[string]interface{}{}
	rows, err = stmt.Paginate(page, limit).Rows()
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	scanData := make([]interface{}, len(columns), len(columns))
	for i := 0; i < len(columns); i++ {
		scanData[i] = new(string)
	}
	for rows.Next() {
		err = rows.Scan(scanData...)
		if err != nil {
			return nil, err
		}
		m := make(map[string]interface{}, len(columns))
		for k, v := range columns {
			m[v] = *scanData[k].(*string)
		}
		p.Data = append(p.Data, m)
	}
	return p, nil
}

//将Paginator数据封装到结构中
func (p *Paginator) Load(value interface{}) error {
	if p == nil {
		return errors.New("Paginator is nil")
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	//封装分页信息
	pageRes := v.FieldByName("PageRes")
	if pageRes.IsValid() {
		if pageRes.Kind() == reflect.Struct {
			page := pageRes.FieldByName("Page")
			if page.IsValid() {
				page.SetUint(p.Page)
			}
			limit := pageRes.FieldByName("Limit")
			if limit.IsValid() {
				limit.SetUint(p.Limit)
			}
			total := pageRes.FieldByName("Total")
			if total.IsValid() {
				total.SetUint(p.Total)
			}
		}
	}

	//封装数据
	if len(p.Data) == 0 {
		return nil
	}
	data := v.FieldByName("Data")
	if data.IsValid() {
		if data.Kind() == reflect.Slice {
			dataPtr := reflect.New(data.Type())
			err := scan.MapsToStructs(p.Data, dataPtr.Interface())
			if err != nil {
				return err
			}
			data.Set(dataPtr.Elem())
			return nil
		}
	}
	return errors.New("no field Data or Data isn't slice")
}

//分页请求体
type PageReq struct {
	//当前页码
	Page uint64 `json:"page"`
	//每页条数
	Limit uint64 `json:"limit"`
}

//分页响应体
type PageRes struct {
	//当前页码
	Page uint64 `json:"page"`
	//每页条数
	Limit uint64 `json:"limit"`
	//总条数
	Total uint64 `json:"total"`
}
