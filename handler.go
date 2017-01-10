/*
   @brief 会话，事物，以及数据库句柄
*/
package misakago

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"gopkg.in/mgo.v2"
)

//数据库 handler
type DbHandler struct {
	Session *mgo.Session  //会话
	DB      *mgo.Database //数据库
}

//会话Handler 是包含一些请求上下文的结构体.
type Handler struct {
	DbHandler                     //数据库句柄
	http.ResponseWriter           //http 响应句柄
	*http.Request                 //http 请求句柄
	StartTime           time.Time //接受请求时间
}

// NewHandler返回含有请求上下文的数据库Handler.
func NewDbHandler() *DbHandler {
	//create session
	session, err := mgo.Dial(Config.DB)
	if err != nil {
		panic(err)
	}
	//set session mode
	session.SetMode(mgo.Monotonic, true)
	return &DbHandler{
		Session: session,
		DB:      session.DB(Config.DbName),
	}
}

// NewHandler返回含有请求上下文的Handler.
func NewHandler(w http.ResponseWriter, r *http.Request) *Handler {
	return &Handler{
		DbHandler:      *NewDbHandler(),
		ResponseWriter: w,
		Request:        r,
		StartTime:      time.Now(),
	}
}

// 只用file作模板的简易渲染
func (handler *Handler) Render(file string, datas ...map[string]interface{}) {
	var data = make(map[string]interface{})
	if len(datas) == 1 {
		data = datas[0]
	} else if len(datas) != 0 {
		panic("不能传入超过多个data map")
	}
	tpl, err := template.ParseFiles(file)
	if err != nil {
		panic(err)
	}
	tpl.Execute(handler.ResponseWriter, data)
}

// param 返回在url中name的值.
func (handler *Handler) Param(name string) string {
	return mux.Vars(handler.Request)[name]
}

// 重定向.
func (handler *Handler) Redirect(urlStr string, code int) {
	http.Redirect(handler.ResponseWriter, handler.Request, urlStr, code)
}

// 返回json数据.
func (handler *Handler) RenderJson(data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	handler.ResponseWriter.Header().Set("Content-Type", "application/json")
	handler.ResponseWriter.Write(b)
}

// 返回文本数据
func (handler *Handler) RenderText(text string) {
	handler.ResponseWriter.Write([]byte(text))
}

//没找到
func (handler *Handler) NotFound() {
	http.NotFound(handler.ResponseWriter, handler.Request)
}

// 渲染模板，并放入一些模板常用变量
func (handler *Handler) RenderTemplate(file string, datas ...map[string]interface{}) {
	var data = make(map[string]interface{})
	data["now"] = time.Now()

	page := parseTemplate(file, data)
	handler.ResponseWriter.Write(page)
}

// Redirect是重定向的简便方法.
func (handler *Handler) RedirectNotfound(urlStr string) {
	http.Redirect(handler.ResponseWriter, handler.Request, urlStr, http.StatusFound)
}

// 获取requestUrl 匹配上的后缀 /topic/p/3
func (handler *Handler) ParseUrlParam() (map[string]string, error) {
	vars := mux.Vars(handler.Request)
	return vars, nil
}

// 获取requestUrl 的标准后缀(as ?x=1&y=3)
func (handler *Handler) ParseExtParam() (url.Values, error) {
	queryForm, err := url.ParseQuery(handler.Request.URL.RawQuery)
	return queryForm, err
}

//文件保存句柄
func fileHandler(w http.ResponseWriter, req *http.Request) {
	url := req.Method + " " + req.URL.Path
	logger.Println(url)
	filePath := req.URL.Path[1:]
	http.ServeFile(w, req, filePath)
}
