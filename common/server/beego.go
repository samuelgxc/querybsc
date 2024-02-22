package server

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"shangchenggo/conf"
	_ "shangchenggo/routers/admin"
	_ "shangchenggo/routers/api"
)

//启动 beego 服务
func StartBeego() {
	beego.Info("--- 启动 beego 服务 ---", "端口:", beego.BConfig.Listen.HTTPPort)
	/* 启动 beego pprof */
	StartBeegoPprof()
	/* 定义模板函数 */
	StartBeegoTplfunc()

	/* 注册 beego 静态文件路由 */
	beego.SetStaticPath("/upload", conf.GetUploadPath())

	/* 注册 beego 过滤器 */
	// beego 配置跨域请求共享CORS过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true, //允许访问所有域
		AllowMethods:    []string{"GET", "POST"},
		AllowHeaders:    []string{"Content-Type", "Accept-Language"}, //header的类型
	}))

	/* 启动 beego 服务 */
	beego.Run()
}
