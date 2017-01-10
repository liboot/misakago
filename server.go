/*
读取配置文件,设置URL,启动服务器
*/

package misakago

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

//初始化日志句柄
var (
	logger = log.New(os.Stdout, "[misakago]:", log.LstdFlags)
)

func LogInfo(format string, v ...interface{}) {
	logger.Println(format, v)
}

//装饰route,构建回调句柄
func handlerFun(route Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//创建回话句柄
		handler := NewHandler(w, r)

		//完成后需要关闭会话
		defer handler.Session.Close()

		url := r.Method + " " + r.URL.Path
		if r.URL.RawQuery != "" {
			url += "?" + r.URL.RawQuery
		}
		logger.Println(url)
		route.HandlerFunc(handler)
	}
}

type MisakaServer struct {
	//配置文件路径
	ConfigPath string
	//路由表
	Routes []Route
}

func (server *MisakaServer) SetConfPath(path string) {
	server.ConfigPath = path
}

func (server *MisakaServer) AddRoute(url string, permission PerType, handler HandlerFunc) {
	route := Route{url, permission, handler}
	server.Routes = append(server.Routes, route)
}

//开启服务
func (server *MisakaServer) StartServer() {

	//解析配置文件
	parseJsonFile(server.ConfigPath, &Config)

	//注册路由表
	r := mux.NewRouter()
	for _, route := range server.Routes {
		logger.Println("register route url success:", route.URL)
		r.HandleFunc(route.URL, handlerFun(route))
	}

	//静态文件路由表
	r.PathPrefix("/static/").HandlerFunc(fileHandler)
	http.Handle("/", r)

	//开始对外服务
	logger.Println("Server start wait to run:", Config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", Config.Port), nil)
	if err != nil {
		logger.Fatal(err)
	}
}
