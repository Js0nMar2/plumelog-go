package log

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"os"
	"plumelog/config"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	debug         *log.Logger
	info          *log.Logger
	warn          *log.Logger
	error         *log.Logger
	dayChangeLock sync.RWMutex
)

const (
	debugLevel = iota //iota=0
	infoLevel
	warnLevel
	errorLevel
)

func init() {
	dayChangeLock = sync.RWMutex{}
	createLogFile()
	go logJob()
}

func createLogFile() {
	dayChangeLock.Lock()
	defer dayChangeLock.Unlock()
	now := time.Now()
	postFix := now.Format("20060102")
	logFile := "log/plume_log_" + postFix + ".log"
	logOut, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	} else {
		writer := io.MultiWriter(os.Stdout, logOut)
		debug = log.New(writer, "[DEBUG] ", log.LstdFlags|log.Lmicroseconds)
		info = log.New(writer, "[INFO] ", log.LstdFlags|log.Lmicroseconds)
		warn = log.New(writer, "[WARN] ", log.LstdFlags|log.Lmicroseconds)
		error = log.New(writer, "[ERROR] ", log.LstdFlags|log.Lmicroseconds)
	}
}

func Debug(format string, v ...any) {
	if config.Conf.Level <= debugLevel {
		debug.Printf(getLineNo()+format, v...)
	}
}

func Info(format string, v ...any) {
	if config.Conf.Level <= infoLevel {
		info.Printf(getLineNo()+format, v...)
	}
}

func Warn(format string, v ...any) {
	if config.Conf.Level <= warnLevel {
		warn.Printf(getLineNo()+format, v...)
	}
}

func Error(format string, v ...any) {
	if config.Conf.Level <= errorLevel {
		error.Printf(getLineNo()+format, v...)
	}
}

func getLineNo() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		split := strings.Split(file, "/")
		file = split[len(split)-1]
		fileLine := file + ":" + strconv.Itoa(line) + " "
		return fileLine
	}
	return ""
}

// logJob 定时操作日志
func logJob() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("@daily", func() {
		Info("执行log定时任务。。。")
		now := time.Now()
		createLogFile()
		//关闭昨天的日志
		closeYesterdayLogFile := fmt.Sprintf("plume_log_%s.log", now.Add(-24*time.Hour).Format("20060102"))
		file, _ := os.Open(closeYesterdayLogFile)
		file.Sync()
		file.Close()
		Info("关闭日志: %s", closeYesterdayLogFile)
		// 删除n天前的日志
		removeLogFile := fmt.Sprintf("plume_log_%s.log", now.Add(time.Duration(config.Conf.Log.KeepDays)*-24*time.Hour).Format("20060102"))
		// 日志是否存在
		_, err := os.Open(removeLogFile)
		if err != nil {
			Error(err.Error())
			return
		}
		go func() {
			// 设置for select 的原因是文件虽然被关闭了，但文件所占的process还在进行中，每10秒轮询一次，执行删除操作，确保文件有被删除
		loop:
			for {
				select {
				case <-time.After(10 * time.Second):
					removeErr := os.Remove(removeLogFile)
					if removeErr != nil {
						Error(removeErr.Error())
					} else {
						Info("删除日志成功：%s", removeLogFile)
						break loop
					}
				}
			}
		}()
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

// 输出有颜色的字体
//func ColorPrint4Window(s, t string) {
//	switch t {
//	case "DEBUG":
//		proc.Call(uintptr(syscall.Stdout), uintptr(FontColor.light_cyan))
//		Debug(s)
//	case "INFO":
//		proc.Call(uintptr(syscall.Stdout), uintptr(FontColor.green))
//		Info(s)
//	case "WARN":
//		proc.Call(uintptr(syscall.Stdout), uintptr(FontColor.light_yellow))
//		Warn(s)
//	case "ERROR":
//		proc.Call(uintptr(syscall.Stdout), uintptr(FontColor.red))
//		Error(s)
//	default:
//		proc.Call(uintptr(syscall.Stdout), uintptr(FontColor.black))
//		Info(s)
//	}
//}
