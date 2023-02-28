package conf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const LoginSessionName = "LoginSessionName"

const CaptchaSessionName = "__captcha__"

const RegexpEmail = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

//允许用户名中出现点号
const RegexpAccount = `^[a-zA-Z][a-zA-Z0-9\.-]{2,50}$`

// PageSize 默认分页条数.
const PageSize = 10

//用户权限
const (
	// 超级管理员.
	MemberSuperRole SystemRole = iota
	//普通管理员.
	MemberAdminRole
	//普通用户.
	MemberGeneralRole
)
const (
	MemberleeRole SystemRole = iota
)

//系统角色
type SystemRole int

const (
	// 创始人.
	BookFounder BookRole = iota
	//管理者
	BookAdmin
	//编辑者.
	BookEditor
	//观察者
	BookObserver
)

//项目角色
type BookRole int

const (
	LoggerOperate   = "operate"
	LoggerSystem    = "system"
	LoggerException = "exception"
	LoggerDocument  = "document"
)
const (
	//本地账户校验
	AuthMethodLocal = "local"
	//LDAP用户校验
	AuthMethodLDAP = "ldap"
)

var (
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
)

var (
	ConfigurationFile = "./conf/app.conf"
	WorkingDirectory  = "./"
	LogFile           = "./runtime/logs"
	BaseUrl           = ""
	AutoLoadDelay     = 0
)

func init() {
	viper.SetConfigFile(ConfigurationFile)
	viper.SetConfigType("toml")
	viper.AddConfigPath(WorkingDirectory)
	err := viper.ReadInConfig()
	if err != nil {
		panic("error to load config file: " + err.Error())
	}
	if viper.GetBool("load_success") {
		log.Println("load config file success")
	} else {
		log.Println("load config file failed")
	}
}

func URLForWithCdnImage(p string) string {
	if strings.HasPrefix(p, "http//") || strings.HasPrefix(p, "https//") {
		return p
	}
	cdn := GetConfig("cdnimg", "").(string)
	if cdn == "" {
		baseUrl := GetConfig("baseurl", "/").(string)

		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseUrl, "/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseUrl, "/") {
			return baseUrl + "/" + p
		}
		return baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

//设置默认值，若配置文件中无值，则使用默认值
func GetConfig(key string, d_value interface{}) interface{} {
	viper.SetDefault(key, d_value)
	return viper.Get(key)
}

//app_key
func GetAppkey() string {
	return GetConfig("app_key", "leedoc").(string)
}

//GEt databaseprefix
func GetDatabasePrefix() string {
	return GetConfig("db_prefix", "lee_").(string)
}

//获取默认头像 avatar
func GetDefaultAvatar() string {
	return GetConfig("avatar", "static/images/headimgurl.jpg").(string)
}

//获取阅读令牌长度.
func GetTokenSize() int {
	return GetConfig("token_size", 12).(int)
}

//获取默认文档封面.
func GetDefaultCover() string {

	return URLForWithCdnImage(GetConfig("cover", "/static/images/book.jpg").(string))
}

//获取允许的文件类型.
func GetUploadFileExt() []string {
	ext := GetConfig("upload_file_ext", "png|jpg|jpeg|gif|txt|doc|docx|pdf").(string)

	temp := strings.Split(ext, "|")

	exts := make([]string, len(temp))

	i := 0
	for _, item := range temp {
		if item != "" {
			exts[i] = item
			i++
		}
	}
	return exts
}

//获取允许上传的文件大小.
func GetUploadFileSize() int64 {
	size := GetConfig("upload_file_size", "10").(string)
	if strings.HasPrefix(size, "MB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024 * 1024
		}
	}
	if strings.HasSuffix(size, "MB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024 * 1024
		}
	}
	if strings.HasSuffix(size, "GB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024 * 1024 * 1024
		}
	}
	if strings.HasSuffix(size, "KB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024
		}
	}
	if s, e := strconv.ParseInt(size, 10, 64); e == nil {
		return s * 1024
	}
	return 0
}

//是否启用导出
func GetEnableExport() bool {
	return GetConfig("enable_export", true).(bool)
}

//同一项目导出线程的并发数
func GetExportProcessNum() int {
	exportProcessNum := GetConfig("export_process_num", 1).(int)

	if exportProcessNum <= 0 || exportProcessNum > 4 {
		exportProcessNum = 1
	}
	return exportProcessNum
}

//导出项目队列的并发数量
func GetExportLimitNum() int {
	exportLimitNum := GetConfig("export_limit_num", 1).(int)

	if exportLimitNum < 0 {
		exportLimitNum = 1
	}
	return exportLimitNum
}

//等待导出队列的长度
func GetExportQueueLimitNum() int {
	exportQueueLimitNum := GetConfig("export_queue_limit_num", 10).(int)

	if exportQueueLimitNum <= 0 {
		exportQueueLimitNum = 100
	}
	return exportQueueLimitNum
}

//默认导出项目的缓存目录
func GetExportOutputPath() string {
	exportOutputPath := filepath.Join(GetConfig("export_output_path", filepath.Join(WorkingDirectory, "cache")).(string), "books")

	return exportOutputPath
}

//判断是否是允许商城的文件类型.
func IsAllowUploadFileExt(ext string) bool {

	if strings.HasPrefix(ext, ".") {
		ext = string(ext[1:])
	}
	exts := GetUploadFileExt()

	for _, item := range exts {
		if item == "*" {
			return true
		}
		if strings.EqualFold(item, ext) {
			return true
		}
	}
	return false
}

func URLForWithCdnCss(p string, v ...string) string {
	cdn := GetConfig("cdncss", "").(string)
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	filePath := WorkingDir(p)

	if f, err := os.Stat(filePath); err == nil && !strings.Contains(p, "?") && len(v) > 0 && v[0] == "version" {
		p = p + fmt.Sprintf("?v=%s", f.ModTime().Format("20060102150405"))
	}
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := GetConfig("baseurl", "/").(string)

		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseUrl, "/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseUrl, "/") {
			return baseUrl + "/" + p
		}
		return baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

func URLForWithCdnJs(p string, v ...string) string {
	cdn := GetConfig("cdnjs", "").(string)
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}

	filePath := WorkingDir(p)

	if f, err := os.Stat(filePath); err == nil && !strings.Contains(p, "?") && len(v) > 0 && v[0] == "version" {
		p = p + fmt.Sprintf("?v=%s", f.ModTime().Format("20060102150405"))
	}

	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := GetConfig("baseurl", "/").(string)

		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseUrl, "/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseUrl, "/") {
			return baseUrl + "/" + p
		}
		return baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}
func WorkingDir(elem ...string) string {

	elems := append([]string{WorkingDirectory}, elem...)

	return filepath.Join(elems...)
}

func init() {
	if p, err := filepath.Abs("./conf/app.conf"); err == nil {
		ConfigurationFile = p
	}
	if p, err := filepath.Abs("./"); err == nil {
		WorkingDirectory = p
	}
	if p, err := filepath.Abs("./runtime/logs"); err == nil {
		LogFile = p
	}
}
