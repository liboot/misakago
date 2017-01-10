/*
   @brief URL和Handler的Mapping
*/

package misakago

import ()

// HandlerFunc 是自定义的请求处理函数,接受*Handler作为参数.
type HandlerFunc func(*Handler)

// Route 是代表对应请求的路由规则以及权限的结构体.
type Route struct {

	//route url 路径
	URL string

	//route 权限标识
	Permission PerType

	//route 的回调函数
	HandlerFunc HandlerFunc
}
