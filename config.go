package misakago

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"os"
	"runtime"
)

type ConfigStruct struct {
	Host                 string `json:"host"`
	Port                 int    `json:"port"`
	DB                   string `json:"db"`
	DbName               string `json:"dbname"`
	LogPath              string `json:"log_path"`
	TemplateDir          string `json:"template_dir"`
	StaticDir            string `json:"static_dir"`
	FileDir              string `json:"file_dir"`
	CookieSecret         string `json:"cookie_secret"`
	Superusers           string `json:"superusers"`
	TimeZoneOffset       int64  `json:"time_zone_offset"`
	StaticFileVersion    int    `json:"static_file_version"`
	GoGetPath            string `json:"go_get_path"`
	PackagesDownloadPath string `json:"packages_download_path"`
	CookieSecure         bool   `json:"cookie_secure"`
	GoDownloadPath       string `json:"go_download_path"`
}

var (
	Config        ConfigStruct
	analyticsCode template.HTML // 网站统计分析代码
	shareCode     template.HTML // 分享代码
	goVersion     = runtime.Version()
)

func GetConfig() ConfigStruct {
	return Config
}

func parseJsonFile(path string, v interface{}) {
	file, err := os.Open(path)
	if err != nil {
		logger.Fatal("配置文件读取失败:", err)
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	err = dec.Decode(v)
	if err != nil {
		logger.Fatal("配置文件解析失败:", err)
	}
}

func getDefaultCode(path string) (code template.HTML) {
	if path != "" {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			logger.Fatal("文件 " + path + " 没有找到")
		}
		code = template.HTML(string(content))
	}
	return
}
