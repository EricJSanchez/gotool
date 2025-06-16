package sys

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type lineHook struct {
	Field string
	// skip为遍历调用栈开始的索引位置
	Skip   int
	levels []logrus.Level
}

var stdoutLog *logrus.Logger

func InitLog() {
	stdoutLog = logrus.New()
	// 添加错误级别高于等于Error的日志HOOK
	stdoutLog.AddHook(new(lineHook))
	// 设置格式为JSON
	stdoutLog.SetFormatter(&logrus.JSONFormatter{})

	//调试写入文件
	stdoutLog.AddHook(new(LogHook))

}

func Log() *logrus.Logger {
	return stdoutLog
}

// Levels implement levels
func (hook lineHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Fire implement fire
func (hook lineHook) Fire(entry *logrus.Entry) error {
	ws := findCaller(hook.Skip, true)
	entry.Data["funcLine"] = ws
	entry.Data["message"] = entry.Message
	logMsg, _ := json.Marshal(entry.Data)
	// 发送错误消息
	if Cfg("app").GetString("ErrNoticeRdsKey") != "" {
		_ = Redis().LPush(context.Background(), Cfg("app").GetString("ErrNoticeRdsKey"), "【"+Cfg("app").GetString("service_name")+"】"+string(logMsg)).Err()
	}
	return nil
}

func findCaller(skip int, lineNew bool) string {
	file := ""
	line := 0
	//var pc uintptr
	rs := ""
	// 遍历调用栈的最大索引为第11层.
	// 参考 1. https://blog.csdn.net/wslyk606/article/details/81670713
	// 参考 2. https://blog.csdn.net/qq_39787367/article/details/109609511
	for i := 0; i < 11; i++ {
		file, line, _ = getCaller(skip + i)
		// 过滤掉所有logrus包，即可得到生成代码信息
		if strings.HasPrefix(file, "logrus") {
			continue
		} else if strings.HasPrefix(file, "log") {
			continue
		} else if strings.HasPrefix(file, "sys") {
			continue
		} else if strings.HasPrefix(file, "runtime") {
			continue
		} else if strings.HasPrefix(file, "datasource") {
			continue
		} else if strings.HasPrefix(file, "task") {
			continue
		} else {
			rTmp := fmt.Sprintf("%s:%d", file, line)
			if lineNew {
				rs = rs + " | file: " + rTmp
			} else {
				rs = rs + " file: " + rTmp
			}
		}
	}
	return rs
}

func getCaller(skip int) (string, int, uintptr) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0, pc
	}
	n := 0

	// 获取包名
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line, pc
}

type LogHook struct {
}

func (hook *LogHook) Fire(entry *logrus.Entry) error {
	pathInfo := Cfg("app").GetString("log_path") + Cfg("app").GetString("service_name")
	l := Logger{pathInfo: pathInfo}
	// 将entry.Data合并到日志消息
	entry.Data["message"] = entry.Message
	logMsg, _ := json.Marshal(entry.Data)
	l.writeToLog(entry.Level, string(logMsg))
	return nil
}

func (hook *LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Logger 按照日期存储日志信息
type Logger struct {
	filename string     // 文件名称
	pathInfo string     // 存储路径
	file     *os.File   //文件指针
	time     *time.Time // 文件日期
}

// 向日志中追加内容
func (this *Logger) writeToLog(l logrus.Level, msg string) {
	this.newLogFile()
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		_, err := fmt.Fprintln(this.file, "[", time.Now().Format("2006-01-02 15:04:05"), "]", "[ERROR]", "[", file, ":", line, "]", "runtime.Caller() fail")
		if err != nil {
			return
		}
		return
	}
	info := findCaller(1, false)
	// 日志信息写入文件中
	_, err := fmt.Fprintln(this.file, "[", time.Now().Format("2006-01-02 15:04:05"), "]", "[", l, "]", "[", info, "]", msg)
	if err != nil {
		return
	}
}

// 若当前是新的一天，则需要创建新文件，同时更新文件信息
func (this *Logger) newLogFile() *Logger {
	// 获取当前日期
	now := time.Now()
	filename := "print_" + now.Format("20060102") + ".log"
	// 获取日志路径
	road := this.pathInfo
	newPath := path.Join(road, filename)
	// 创建新的文件 以当前年月日命名
	_, errExt := os.Stat(filepath.Dir(newPath))
	if os.IsNotExist(errExt) {
		errTmp := os.MkdirAll(filepath.Dir(newPath), 0755)
		if errTmp != nil {
			fmt.Println("MkdirAll Err:", errTmp)
		}
	}
	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// 将err作为日志输出在文件中
		_, f, line, _ := runtime.Caller(0)
		_, err := fmt.Fprintln(this.file, "[", time.Now().Format("2006-01-02 15:04:05"), "]", "[ERROR]", "[", f, ":", line, "]", "create log file failed!")
		if err != nil {
			return nil
		}
		return this
	}
	// 更新结构体数据
	this.filename = filename
	this.file = file
	this.time = &now
	return this
}
