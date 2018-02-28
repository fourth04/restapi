package middlewares

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)

const deleteFileOnExit = false

// ConfigPlus is a extension of logger.Config, added logPath
type ConfigPlus struct {
	LogPath string
}

// get a filename based on the date, file logs works that way the most times
// but these are just a sugar.
func todayFilename() string {
	today := time.Now().Format("20060602")
	return today + ".txt"
}

func newLogFile(logPath string) *os.File {
	filename := todayFilename()
	// open an output file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(logPath+"/"+filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}

var excludeExtensions = [...]string{
	".js",
	".css",
	".jpg",
	".png",
	".ico",
	".svg",
}

func myColumnize(nowFormatted string, latency time.Duration, status, ip, method, path string, message interface{}) string {

	// Time | Status | Latency | IP | Method | Path
	line := fmt.Sprintf("%s|%v|%v|%s|%s|%s", nowFormatted, status, latency, ip, method, path)
	if message != nil {
		line += fmt.Sprintf("|%v", message)
	}

	output := line + "\n"
	return output
}

// NewRequestLogger 用于返回一个自定义的Logger中间件
func NewRequestLogger(c *logger.Config, cp *ConfigPlus) (h iris.Handler, close func() error) {
	close = func() error { return nil }

	logFile := newLogFile(cp.LogPath)
	close = func() error {
		err := logFile.Close()
		if deleteFileOnExit {
			err = os.Remove(logFile.Name())
		}
		return err
	}

	c.LogFunc = func(now time.Time, latency time.Duration, status, ip, method, path string, message interface{}) {
		output := myColumnize(now.Format("2006/01/02-15:04:05"), latency, status, ip, method, path, message)
		logFile.Write([]byte(output))
	}

	// we don't want to use the logger
	// to log requests to assets and etc
	c.AddSkipper(func(ctx iris.Context) bool {
		path := ctx.Path()
		for _, ext := range excludeExtensions {
			if strings.HasSuffix(path, ext) {
				return true
			}
		}
		return false
	})

	h = logger.New(*c)

	return
}
