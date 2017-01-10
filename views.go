/*
   view 层的一些辅助方法
*/

package misakago

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/pborman/uuid"
	qiniuIo "github.com/qiniu/api.v6/io"
	"github.com/qiniu/api.v6/rs"
)

var (
	store       *sessions.CookieStore
	fileVersion map[string]string = make(map[string]string) // {path: version}
	usersJson   []byte
)

// 获取链接的页码，默认"?p=1"这种类型
func Page(r *http.Request) (int, error) {
	p := r.FormValue("p")
	page := 1

	if p != "" {
		var err error
		page, err = strconv.Atoi(p)

		if err != nil {
			return 0, err
		}
	}

	return page, nil
}

// 检查一个string元素是否在数组里面
func stringInArray(a []string, x string) bool {
	sort.Strings(a)
	index := sort.SearchStrings(a, x)

	if index == 0 {
		if a[0] == x {
			return true
		}

		return false
	} else if index > len(a)-1 {
		return false
	}

	return true
}

func StaticHandler(templateFile string) HandlerFunc {
	return func(handler *Handler) {
		handler.RenderTemplate(templateFile, map[string]interface{}{})
	}
}

func getPage(r *http.Request) (page int, err error) {
	p := r.FormValue("p")
	page = 1

	if p != "" {
		page, err = strconv.Atoi(p)

		if err != nil {
			return
		}
	}

	return
}

// 编辑器上传图片，接收后上传到七牛
func uploadImageHandler(handler *Handler) {
	file, header, err := handler.Request.FormFile("editormd-image-file")
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()

	// 检查是否是jpg或png文件
	uploadFileType := header.Header["Content-Type"][0]

	filenameExtension := ""
	if uploadFileType == "image/jpeg" {
		filenameExtension = ".jpg"
	} else if uploadFileType == "image/png" {
		filenameExtension = ".png"
	} else if uploadFileType == "image/gif" {
		filenameExtension = ".gif"
	}

	if filenameExtension == "" {
		handler.RenderJson(map[string]interface{}{
			"success": 0,
			"message": "不支持的文件格式，请上传 jpg/png/gif 图片",
		})
		return
	}

	// 上传到七牛
	// 文件名：32位uuid+后缀组成
	filename := strings.Replace(uuid.NewUUID().String(), "-", "", -1) + filenameExtension
	key := "upload/image/" + filename

	ret := new(qiniuIo.PutRet)

	var policy = rs.PutPolicy{
		Scope: "gopher",
	}

	err = qiniuIo.Put(
		nil,
		ret,
		policy.Token(nil),
		key,
		file,
		nil,
	)

	if err != nil {
		panic(err)

		handler.RenderJson(map[string]interface{}{
			"success": 0,
			"message": "图片上传到七牛失败",
		})

		return
	}

	handler.RenderJson(map[string]interface{}{
		"success": 1,
		"url":     "http://77fkk5.com1.z0.glb.clouddn.com/" + key,
	})
}
