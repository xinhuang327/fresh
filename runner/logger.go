package runner

import (
	"fmt"
	logPkg "log"
	"time"

	"github.com/mattn/go-colorable"
)

type logFunc func(string, ...interface{})

var logger = logPkg.New(colorable.NewColorableStderr(), "", 0)

func newLogFunc(prefix string) func(string, ...interface{}) {
	isApp := prefix == "app"
	color, clear := "", ""
	if settings["colors"] == "1" {
		color = fmt.Sprintf("\033[%sm", logColor(prefix))
		clear = fmt.Sprintf("\033[%sm", colors["reset"])
	}
	prefix = fmt.Sprintf("%-5s", prefix)

	if isApp {
		return func(format string, v ...interface{}) {
			format = fmt.Sprintf("%s%s", color, format)
			//logger.Printf(format, v...)
			logger.Print(format)
		}
	} else {
		return func(format string, v ...interface{}) {
			now := time.Now()
			timeString := fmt.Sprintf("%d:%d:%02d", now.Hour(), now.Minute(), now.Second())
			format = fmt.Sprintf("%s%s %s |%s %s", color, timeString, prefix, clear, format)
			logger.Printf(format, v...)
		}
	}
}

func fatal(err error) {
	logger.Fatal(err)
}

type appLogWriter struct{}

func (a appLogWriter) Write(p []byte) (n int, err error) {
	appLog(string(p))

	return len(p), nil
}
