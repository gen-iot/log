package log

import (
	"fmt"
	"gitee.com/gen-iot/std"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

/**
 * Created by suzhen on 2018-06-15 .
 */

const kDefaultLogTimeFormat = "2006-01-02T15:04:05.999999999-07:00"
const kLoggerTag = "[Logger]"
const kDefaultLogDir = "logs"

var DefaultConfig = Config{
	Level:     logger.InfoLevel,
	LogToFile: false,
	LogsDir:   kDefaultLogDir,
}

type Config struct {
	Level     logger.Level `json:"level" validate:"required"`
	LogToFile bool         `json:"logToFile"`
	LogsDir   string       `json:"logsDir"`
	Timestamp bool         `json:"timestamp"`
}

var gConf = DefaultConfig

type Logger interface {
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

type LoggerProxy struct {
	PrintlnCb func(v ...interface{})
	PrintfCb  func(format string, v ...interface{})
}

func (l *LoggerProxy) Print(v ...interface{}) {
	l.Println(v...)
}
func (l *LoggerProxy) Println(v ...interface{}) {
	l.PrintlnCb(v...)
}
func (l *LoggerProxy) Printf(format string, v ...interface{}) {
	l.PrintfCb(format, v...)
}

var STD Logger = &LoggerProxy{
	PrintlnCb: func(v ...interface{}) {
		fmt.Println("STD ["+time.Now().Format(kDefaultLogTimeFormat)+"]", v)
	},
	PrintfCb: func(format string, v ...interface{}) {
		fmt.Printf(format, v...)
	},
}

var DEBUG Logger = &EmptyLogger{}

var INFO Logger = &EmptyLogger{}

var WARN Logger = &EmptyLogger{}

var ERROR Logger = &EmptyLogger{}

//noinspection GoUnusedGlobalVariable
var FATAL Logger = &EmptyLogger{}
var PANIC Logger = &EmptyLogger{}

func Init() {
	InitWithConfig(DefaultConfig)
}

func InitWithConfig(config Config) {
	if err := std.ValidateStruct(gConf); err != nil {
		panic(errors.WithMessage(err, "logger配置不正确"))
	}
	STD.Println(kLoggerTag, "logger(level=", config.Level, ",toFile=", config.LogToFile, ",logsDir=", config.LogsDir, ") init ... ")
	logger.SetLevel(config.Level)
	logger.SetFormatter(&logger.TextFormatter{
		EnvironmentOverrideColors: true,
		ForceColors:               true,
		FullTimestamp:             config.Timestamp,
		DisableTimestamp:          !config.Timestamp,
		TimestampFormat:           kDefaultLogTimeFormat,
	})
	logger.SetOutput(os.Stdout)
	if config.LogToFile {
		dir := config.LogsDir
		if len(dir) == 0 {
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			dir = filepath.Join(wd, kDefaultLogDir)
		}
		logFile := filepath.Join(dir, "logs.log")
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			STD.Println(kLoggerTag, "redirect log to file error :", err)
		}
		logger.SetOutput(file)
	}
	DEBUG = &LoggerProxy{PrintlnCb: logger.Debug, PrintfCb: logger.Debugf}
	INFO = &LoggerProxy{PrintlnCb: logger.Info, PrintfCb: logger.Infof}
	WARN = &LoggerProxy{PrintlnCb: logger.Warn, PrintfCb: logger.Warnf}
	ERROR = &LoggerProxy{PrintlnCb: logger.Error, PrintfCb: logger.Errorf}
	FATAL = &LoggerProxy{PrintlnCb: logger.Fatal, PrintfCb: logger.Fatalf}
	PANIC = &LoggerProxy{PrintlnCb: logger.Panic, PrintfCb: logger.Panicf}
}
