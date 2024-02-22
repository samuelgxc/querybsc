package conf

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//项目根目录绝对路径
var appPath string

func init() {
	//log.Print("=== 加载 beego 配置文件 ===")
	var beegoConfigFile string
	filePath := filepath.Join("conf", "app.conf")
	workPath, _ := os.Getwd() //workDir
	workPath, ok := confExist(workPath, filePath)
	if ok {
		beegoConfigFile = filepath.Join(workPath, filePath)
		SetAppPath(workPath)
	} else {
		workPath, _ = filepath.Abs(filepath.Dir(os.Args[0])) //outputDir
		workPath, ok = confExist(workPath, filePath)
		if ok {
			beegoConfigFile = filepath.Join(workPath, filePath)
			SetAppPath(workPath)
		} else {
			workPath = GetAppPath()
			workPath, ok = confExist(workPath, filePath)
			if ok {
				beegoConfigFile = filepath.Join(workPath, filePath)
			} else {
				log.Fatal("--- 找不到 beego 配置文件 ---", " filePath: ", filePath)
			}
		}
	}
	err := beego.LoadAppConfig("ini", beegoConfigFile)
	if err != nil {
		log.Fatal("--- 加载 beego 配置文件出错 ---", " beegoConfigFile: ", beegoConfigFile, " err: ", err)
	}

	/* 启动 beego 文件日志 */
	beego.SetLogger(logs.AdapterFile, `{"filename":"`+GetLogFilePath("beegolog")+`"}`)
}

func Init() {}

// 设置项目根目录
func SetAppPath(path string) {
	appPath = path
}

// 获取项目根目录
func GetAppPath() string {
	return appPath
}

// 获取日志存放路径
func GetLogPath() string {
	logPath := beego.AppConfig.String(LOG_PATH)
	if logPath != "" {
		if ok, _ := createDir(logPath, 0755); ok {
			return logPath
		}
	}
	logPath = filepath.Join(GetAppPath(), "log")
	if ok, _ := createDir(logPath, 0755); ok {
		return logPath
	}
	return GetAppPath()
}

// 获取日志文件路径
func GetLogFilePath(logName string) string {
	logPath := GetLogPath()
	if logName == "" {
		logName = beego.AppConfig.DefaultString(LOG_NAME, "applog")
	}
	//日志文件全路径名
	logFile := filepath.Join(logPath, fmt.Sprintf("%s.log", logName))
	//logFile := filepath.Join(logPath, fmt.Sprintf("%s_%s.log", logName, time.Now().Format("2006-01-02")))
	return strings.Replace(logFile, "\\", "/", -1)
}

// 获取文件存放绝对路径
func GetUploadPath() string {
	uploadPath := beego.AppConfig.String(UPLOAD_PATH)
	if uploadPath != "" {
		if ok, _ := createDir(uploadPath, 0755); ok {
			return uploadPath
		}
	}
	uploadPath = filepath.Join(GetAppPath(), "upload")
	if ok, _ := createDir(uploadPath, 0755); ok {
		return uploadPath
	}
	return GetAppPath()
}

//获取mysql连接配置
func GetDbConfig() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&allowOldPasswords=%d",
		beego.AppConfig.DefaultString(MYSQL_USERNAME, "gaopeng"),
		beego.AppConfig.DefaultString(MYSQL_PASSWORD, "gaopeng"),
		beego.AppConfig.DefaultString(MYSQL_HOST, "192.168.0.200"),
		beego.AppConfig.DefaultInt(MYSQL_PORT, 3306),
		beego.AppConfig.DefaultString(MYSQL_DATABASE, "gaopeng_demomcc"),
		beego.AppConfig.DefaultString(MYSQL_CHARSET, "utf8mb4"),
		beego.AppConfig.DefaultInt(MYSQL_ALLOWOLDPASSWORDS, 1),
	)
}

//获取mysql连接配置
func GetDbReadConfig() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&allowOldPasswords=%d",
		beego.AppConfig.DefaultString(MYSQL_READ_USERNAME, "gaopeng"),
		beego.AppConfig.DefaultString(MYSQL_READ_PASSWORD, "gaopeng"),
		beego.AppConfig.DefaultString(MYSQL_READ_HOST, "192.168.0.200"),
		beego.AppConfig.DefaultInt(MYSQL_READ_PORT, 3306),
		beego.AppConfig.DefaultString(MYSQL_READ_DATABASE, "gaopeng_demomcc"),
		beego.AppConfig.DefaultString(MYSQL_READ_CHARSET, "utf8mb4"),
		beego.AppConfig.DefaultInt(MYSQL_READ_ALLOWOLDPASSWORDS, 1),
	)
}

func confExist(workPath string, filePath string) (string, bool) {
	workPath = strings.Replace(workPath, `/`, `\`, -1)
	configFile := filepath.Join(workPath, filePath)
	//fmt.Println(configFile)
	if fileExist(configFile) {
		return workPath, true
	} else {
		if workPath == `` || workPath == `\` || strings.HasSuffix(workPath, `:\`) {
			return workPath, false
		} else {
			workPath = workPath[:strings.LastIndex(workPath, `\`)]
			if workPath == `` || strings.HasSuffix(workPath, `:`) {
				workPath = workPath + `\`
			}
			return confExist(workPath, filePath)
		}
	}
}

//判断文件是否存在
func fileExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 递归创建目录
func createDir(path string, perm os.FileMode) (bool, error) {
	_, err := os.Stat(path)
	// 如果返回的错误为nil,说明文件或文件夹存在
	if err == nil {
		return true, nil
	}
	// 如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
	if !os.IsNotExist(err) {
		return true, nil
	}
	// 目录不存在，递归创建
	err = os.MkdirAll(path, perm)
	if err != nil {
		return false, err
	}
	return true, nil
}

// beego 配置文件参数
const (
	LOG_PATH    = "log_path"    //日志文件绝对路径，例如 D:/log ；默认存放在项目根目录下的log文件夹中
	LOG_NAME    = "log_name"    //日志文件名，例如 app ；默认为applog，自动生成扩展名为 applog.log
	UPLOAD_PATH = "upload_path" //上传文件绝对路径，例如 D:/upload ；默认存放在项目根目录下的upload文件夹中

	//mysql主库（写库）配置
	MYSQL_HOST              = "mysql_host"
	MYSQL_PORT              = "mysql_port"
	MYSQL_USERNAME          = "mysql_username"
	MYSQL_PASSWORD          = "mysql_password"
	MYSQL_DATABASE          = "mysql_database"
	MYSQL_CHARSET           = "mysql_charset"
	MYSQL_ALLOWOLDPASSWORDS = "mysql_allowoldpasswords"
	MYSQL_MAXOPENCONNS      = "mysql_maxopenconns"
	MYSQL_MAXIDLECONNS      = "mysql_maxidleconns"
	MYSQL_CONNMAXLIFETIME   = "mysql_connmaxlifetime"
	//mysql从库（读库）配置
	MYSQL_READ_HOST              = "mysql_read_host"
	MYSQL_READ_PORT              = "mysql_read_port"
	MYSQL_READ_USERNAME          = "mysql_read_username"
	MYSQL_READ_PASSWORD          = "mysql_read_password"
	MYSQL_READ_DATABASE          = "mysql_read_database"
	MYSQL_READ_CHARSET           = "mysql_read_charset"
	MYSQL_READ_ALLOWOLDPASSWORDS = "mysql_read_allowoldpasswords"
)

const (
	//redis配置
	REDIS_HOST        = "redis_host"
	REDIS_PORT        = "redis_port"
	REDIS_PASSWORD    = "redis_password"
	REDIS_DB          = "redis_db"
	REDIS_MAXACTIVE   = "redis_maxactive"
	REDIS_MAXIDLE     = "redis_maxidle"
	REDIS_IDLETIMEOUT = "redis_idletimeout"

	//i18n支持的国际化语言类型（用|分割，默认中文zh-CN，对应语言存储在conf/locale_语言.ini中，若未设置对应语言则显示key名）
	LANG_DEFAULT = "lang_default"
	LANG_TYPES   = "lang_types"
)
