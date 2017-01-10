package misakago

import (
	//"fmt"
	"errors"
	"strconv"
	"strings"
)

//辅助工具包
type MisakaUtil struct {
}

//切割空格和换行
func (util *MisakaUtil) splitBlank(s rune) bool {
	if s == '\n' {
		return true
	}
	if s == ' ' {
		return true
	}
	return false
}

//切割字符串的函数
func (util *MisakaUtil) Split(f func(s rune) bool, str string) []string {
	//如果没有传入切割函数则使用默认的
	if f == nil {
		f = util.splitBlank
	}
	return strings.FieldsFunc(str, f)
}

//将数据转化为字符串
func (util *MisakaUtil) ToString(a interface{}) (string, error) {

	if v, p := a.(int); p {
		return strconv.Itoa(v), nil
	}

	if v, p := a.(float64); p {
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	}

	if v, p := a.(float32); p {
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	}

	if v, p := a.(int16); p {
		return strconv.Itoa(int(v)), nil
	}
	if v, p := a.(uint); p {
		return strconv.Itoa(int(v)), nil
	}
	if v, p := a.(int32); p {
		return strconv.Itoa(int(v)), nil
	}
	return "wrong", errors.New("转换失败")
}

/*
func main() {
	var util MisakaUtil
	var str string = "maaa fff               bbq"
	fmt.Println(util.splitBlank(str))
	vec := util.splitBlank(str)
	fmt.Println(vec[0])
}
*/
