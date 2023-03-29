package log

import (
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"os"
	"plumelog/config"
	"time"
)

type Conf struct {
	LogConf `toml:"log"`
}

type LogConf struct {
	ServiceName string `json:"serviceName"`
	Port        int    `json:"port"`
}

var Port int
var Debug *log.Logger
var Info *log.Logger
var Warn *log.Logger
var Error *log.Logger

func init() {
	format := time.Now().Format("20060102")
	logFile, err := os.OpenFile("plume_log_"+format+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	Debug = log.New(multiWriter, "[debug]", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(multiWriter, "[info]", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(multiWriter, "[warn]", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(multiWriter, "[error]", log.Ldate|log.Ltime|log.Lshortfile)
	go LogJob()
}

// LogJob 定时删除日志
func LogJob() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("@daily", func() {
		Info.Println("执行log定时任务。。。")
		format := time.Now().AddDate(0, 0, -config.Conf.KeepDays).Format("20060102")
		err := os.Remove("plume_log_" + format + ".log")
		if err != nil {
			Error.Println(err)
		}
	})
	c.Start()
}

//var (
//	kernel32    = syscall.NewLazyDLL(`kernel32.dll`)
//	proc        = kernel32.NewProc(`SetConsoleTextAttribute`)
//	CloseHandle = kernel32.NewProc(`CloseHandle`)
//	// 给字体颜色对象赋值
//	FontColor = Color{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
//)

//type Color struct {
//	black        int // 黑色
//	blue         int // 蓝色
//	green        int // 绿色
//	cyan         int // 青色
//	red          int // 红色
//	purple       int // 紫色
//	yellow       int // 黄色
//	light_gray   int // 淡灰色（系统默认值）
//	gray         int // 灰色
//	light_blue   int // 亮蓝色
//	light_green  int // 亮绿色
//	light_cyan   int // 亮青色
//	light_red    int // 亮红色
//	light_purple int // 亮紫色
//	light_yellow int // 亮黄色
//	white        int // 白色
//}
//
//// 输出有颜色的字体
//func ColorPrint4Window(s, t string) {
//	switch t {
//	case "DEBUG":
//		proc.Call(uintptr(syscall.Stdout), uintptr(FontColor.light_cyan))
//		Debug.Println(s)
//	case "INFO":
//		proc.Call(uintptr(syscall.Stdout), uintptr(FontColor.green))
//		Info.Println(s)
//	case "WARN":
//		proc.Call(uintptr(syscall.Stdout), uintptr(FontColor.light_yellow))
//		Warn.Println(s)
//	case "ERROR":
//		proc.Call(uintptr(syscall.Stdout), uintptr(FontColor.red))
//		Error.Println(s)
//	default:
//		proc.Call(uintptr(syscall.Stdout), uintptr(FontColor.black))
//		Info.Println(s)
//	}
//}
