package funcmap

import (
	"html/template"
	"leedoc/conf"
	"log"
	"strings"
	"time"

	"github.com/beego/i18n"
	"github.com/gin-gonic/gin"
)

var FuncMap = template.FuncMap{}

func init() {
	// Register the funcmap
	// gin.RegisterFuncMap("cdncss", cdncss)
	FuncMap["cdncss"] = conf.URLForWithCdnCss
	FuncMap["cdnjs"] = conf.URLForWithCdnJs
	FuncMap["cdnimg"] = conf.URLForWithCdnImage
	FuncMap["date_format"] = func(t time.Time, format string) string {
		return t.Local().Format(format)
	}
	FuncMap["urlfor"] = URLFor
	FuncMap["i18n"] = i18n.Tr
	langs := strings.Split("en-us|zh-cn", "|")
	for _, lang := range langs {
		if err := i18n.SetMessage(lang, "conf/lang/"+lang+".ini"); err != nil {
			log.Println(err)
			return
		}
	}
}

//重写生成URL的方法，加上完整的域名
func URLFor(endpoint string, values string) string {
	baseUrl := conf.GetConfig("base_url", "").(string)

	pathurl := urlFor(endpoint, values)
	if strings.HasPrefix(pathurl, "/") && strings.HasSuffix(baseUrl, "/") {
		return baseUrl + pathurl[1:]
	}
	if !strings.HasPrefix(pathurl, "/") && !strings.HasSuffix(baseUrl, "/") {
		return baseUrl + "/" + pathurl
	}
	return baseUrl + pathurl
}
func URLforNotHost(endpoint string, values interface{}) string {
	baseUrl := conf.GetConfig("base_url", "").(string)
	pathurl := urlFor(endpoint, values)

	if strings.HasPrefix(pathurl, "/") && strings.HasSuffix(baseUrl, "/") {
		return baseUrl + pathurl[1:]
	}
	if !strings.HasPrefix(pathurl, "/") && !strings.HasSuffix(baseUrl, "/") {
		return baseUrl + "/" + pathurl
	}
	return baseUrl + pathurl
}

var URLmap map[string]string

func init() {
	URLmap = make(map[string]string)
}
func URL(routes gin.RoutesInfo) {
	for _, v := range routes {
		handle := strings.Split(v.Handler, "/")
		URLmap[handle[len(handle)-1]] = v.Path
	}
}
func urlFor(endpoint string, valuse interface{}) string {
	s := strings.Split(URLmap[endpoint], "/")
	if valuse == nil || valuse == "" {
		if len(s) > 1 {
			return strings.Join(s[:len(s)-1], "/")
		}
		if len(s) == 1 {
			return strings.Join(s, "/")
		}
	} else {
		if len(s) > 1 {
			return strings.Join(s[:len(s)-1], "/") + "/" + valuse.(string)
		}
		if len(s) == 1 {
			return strings.Join(s, "/")
		}
	}
	return ""
}
