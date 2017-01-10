package misakago

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/jimmykuu/webhelpers"
)

//自定义的模板回调函数
var funcMaps = template.FuncMap{
	//@TODO 文件扩展支持
	"html": func(text string) template.HTML {
		return template.HTML(text)
	},
	"loadtimes": func(startTime time.Time) string {
		// 加载时间
		return fmt.Sprintf("%dms", time.Now().Sub(startTime)/1000000)
	},
	"url": func(url string) string {
		// 没有http://或https://开头的增加http://
		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
			return url
		}

		return "http://" + url
	},
	"add": func(a, b int) int {
		// 加法运算
		return a + b
	},
	"formatdate": func(t time.Time) string {
		// 格式化日期
		return t.Format(time.RFC822)
	},
	"formattime": func(t time.Time) string {
		// 格式化时间
		now := time.Now()
		duration := now.Sub(t)
		if duration.Seconds() < 60 {
			return fmt.Sprintf("刚刚")
		} else if duration.Minutes() < 60 {
			return fmt.Sprintf("%.0f 分钟前", duration.Minutes())
		} else if duration.Hours() < 24 {
			return fmt.Sprintf("%.0f 小时前", duration.Hours())
		}

		t = t.Add(time.Hour * time.Duration(Config.TimeZoneOffset))
		return t.Format("2006-01-02 15:04")
	},
	"formatdatetime": func(t time.Time) string {
		// 格式化时间成 2006-01-02 15:04:05
		return t.Add(time.Hour * time.Duration(Config.TimeZoneOffset)).Format("2006-01-02 15:04:05")
	},
	"nl2br": func(text string) template.HTML {
		return template.HTML(strings.Replace(text, "\n", "<br>", -1))
	},
	"truncate": func(text string, length int, indicator string) string {
		return webhelpers.Truncate(text, length, indicator)
	},
	"include": func(filename string, data map[string]interface{}) template.HTML {
		// 加载局部模板，从 templates 中去寻找
		var buf bytes.Buffer
		t, err := template.ParseFiles(Config.TemplateDir + "/" + filename)
		if err != nil {
			panic(err)
		}
		err = t.Execute(&buf, data)
		if err != nil {
			panic(err)
		}
		return template.HTML(buf.Bytes())
	},
}

// 解析模板
func parseTemplate(file string, data map[string]interface{}) []byte {
	var buf bytes.Buffer
	var err error
	t := template.New(file).Funcs(funcMaps)
	templateFile := Config.TemplateDir + "/" + file
	logger.Println(templateFile)
	t, err = t.ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}
	err = t.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// 渲染模板，并放入一些模板常用变量
func renderTemplate(handler Handler, file string, data map[string]interface{}) {
	page := parseTemplate(file, data)
	handler.ResponseWriter.Write(page)
}

//构造json喘作为http请求的回复
func renderJson(handler Handler, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	handler.ResponseWriter.Header().Set("Content-Type", "application/json")
	handler.ResponseWriter.Write(b)
}
