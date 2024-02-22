package server

import (
	"github.com/astaxie/beego"
	"net/http/pprof"
)

//启动 beego pprof
func StartBeegoPprof() {
	//是否开启beego pprof，默认关闭
	if beego.AppConfig.DefaultBool("beego_pprof", true) {
		// 注册beego的pprof路由
		beego.Router(`/debug/pprof/`, &PprofController{})
		beego.Router(`/debug/pprof/:pp([\w]+)`, &PprofController{})
	}
}

type PprofController struct {
	beego.Controller
}

func (c *PprofController) Get() {
	switch c.Ctx.Input.Param(":pp") {
	case "":
		pprof.Index(c.Ctx.ResponseWriter, c.Ctx.Request)
	case "cmdline":
		pprof.Cmdline(c.Ctx.ResponseWriter, c.Ctx.Request)
	case "profile":
		pprof.Profile(c.Ctx.ResponseWriter, c.Ctx.Request)
	case "symbol":
		pprof.Symbol(c.Ctx.ResponseWriter, c.Ctx.Request)
	case "trace":
		pprof.Trace(c.Ctx.ResponseWriter, c.Ctx.Request)
	default:
		pprof.Index(c.Ctx.ResponseWriter, c.Ctx.Request)
	}
	c.Ctx.ResponseWriter.WriteHeader(200)
}
