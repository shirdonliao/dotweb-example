package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/session"
	"net/http"
	"strconv"
	"net/url"
	"strings"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	app.UseRequestLog()
	//设置gzip开关
	//app.HttpServer.SetEnabledGzip(true)

	//设置Session开关
	app.HttpServer.SetEnabledSession(true)

	//设置Session配置
	//runtime mode
	app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())
	//redis mode
	//app.SetSessionConfig(session.NewDefaultRedisConfig("192.168.8.175:6379", ""))

	//设置路由
	InitRoute(app.HttpServer)

	//设置HttpModule
	//InitModule(app)

	//启动 监控服务
	//pprofport := 8081
	//app.SetPProfConfig(true, pprofport)

	app.SetNotFoundHandle(func(context dotweb.Context) {
		context.Response().SetHeader("dotweb-t", "tres")
		context.WriteJsonC(http.StatusNotFound, "Not Found Json")
	})

	//全局容器
	app.Items.Set("gstring", "gvalue")
	app.Items.Set("gint", 1)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func EchoUrl(ctx dotweb.Context) error{
	fmt.Println(ctx.Request().URL)
	fmt.Println(ctx.Request().Url())
	fmt.Println(fmt.Sprintf(fmt.Sprintf("%s", ctx.Request().Url())))

	fmt.Println(url.QueryEscape(ctx.Request().QueryString("pageurl")))

	return ctx.WriteString(ctx.Request().RequestURI)
}

func Index(ctx dotweb.Context) error {
	//ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Println("Index")
	if strings.Index(ctx.Request().Host, ":") >0{
		ctx.WriteString(strings.Split(ctx.Request().Host, ":")[0])
	}else{
		ctx.WriteString(ctx.Request().Host)
	}

	return nil
}

func IndexReg(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString("welcome to dotweb")
	return nil
}

func KeyPost(ctx dotweb.Context) error {
	username1 := ctx.Request().PostString("username")
	username2 := ctx.FormValue("username")
	username3 := ctx.PostFormValue("username")
	ctx.WriteString("username:" + username1 + " - " + username2 + " - " + username3)
	return nil
}

func JsonPost(ctx dotweb.Context) error {
	ctx.WriteString("body:" + string(ctx.Request().PostBody()))
	return nil
}

func DefaultError(ctx dotweb.Context) error {
	panic("my panic error!")
	return nil
}

func Redirect(ctx dotweb.Context) error {
	ctx.Redirect(http.StatusOK, "http://www.baidu.com")
	return nil
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/Page/PageView", EchoUrl)
	server.Router().GET("/", Index)
	server.Router().POST("/keypost", KeyPost)
	server.Router().POST("/jsonpost", JsonPost)
	server.Router().GET("/error", DefaultError)
	server.Router().GET("/redirect", Redirect)
	server.Router().RegisterRoute(dotweb.RouteMethod_GET, "/index", IndexReg)
}

func InitModule(dotserver *dotweb.DotWeb) {
	dotserver.HttpServer.RegisterModule(&dotweb.HttpModule{
		Name:"test change route",
		OnBeginRequest: func(ctx dotweb.Context) {
			fmt.Println("BeginRequest1:", ctx)
		},
		OnEndRequest: func(ctx dotweb.Context) {
			fmt.Println("EndRequest1:", ctx)
		},
	})

	dotserver.HttpServer.RegisterModule(&dotweb.HttpModule{
		OnBeginRequest: func(ctx dotweb.Context) {
			fmt.Println("BeginRequest2:", ctx)
		},
	})
	dotserver.HttpServer.RegisterModule(&dotweb.HttpModule{
		OnEndRequest: func(ctx dotweb.Context) {
			fmt.Println("EndRequest3:", ctx)
		},
	})
}
