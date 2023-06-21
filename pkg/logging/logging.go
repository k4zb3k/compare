package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}

	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var entr *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{entr}
}

func (l *Logger) GetLoggerWithField(key string, value interface{}) Logger {
	return Logger{l.WithField(key, value)}
}

func init() {
	logg := logrus.New()
	logg.SetReportCaller(true)
	logg.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.MkdirAll("/home/k4zb3k/Desktop/compare/logs", 0777)
	if err != nil {
		log.Println(err)
	}

	allFile, err := os.OpenFile("/home/k4zb3k/Desktop/compare/logs/all.logs", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Println(err)
	}

	logg.SetOutput(io.Discard)

	logg.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	logg.SetLevel(logrus.TraceLevel)

	entr = logrus.NewEntry(logg)
}
