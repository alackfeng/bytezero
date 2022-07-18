package utils

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// 创建全局唯一实例日志.
var log = logrus.New()

// Fields logrus.Fields
type Fields map[string]interface{}

// Logger -  返回唯一实例.
func Logger(f map[string]interface{}) *logrus.Entry {
	if f != nil {
		return log.WithFields(logrus.Fields(f))
	}
	return log.WithField("animal", "default")
}

func init() {

	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter)
	log.Formatter.(*logrus.TextFormatter).DisableColors = false
	log.Formatter.(*logrus.TextFormatter).DisableTimestamp = false
	log.Formatter.(*logrus.TextFormatter).FullTimestamp = true
	// log.SetReportCaller(true)

	// log.Formatter = &logrus.JSONFormatter{
	//  DisableTimestamp: true,
	//  // FullTimestamp:    true,
	//  CallerPrettyfier: func(f *runtime.Frame) (string, string) {
	//    s := strings.Split(f.Function, ".")
	//    funcname := s[len(s)-1]
	//    _, filename := path.Split(f.File)
	//    return funcname, filename
	//  },
	// }

	log.Level = logrus.TraceLevel
	// log.Level = logrus.WarnLevel
	log.Out = os.Stdout

	// logFile, err := os.OpenFile("dcs.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	//  log.Out = logFile
	// } else {
	//  log.Info("Failed to log to file, using default stderr")
	// }

	log.WithFields(logrus.Fields{
		"animal": "log",
	}).Info("Init Logger.")
}

// SetLogout -  设置日志输出位置.
func SetLogout(path string) {
    if len(path) < 3 {
        return
    }
	if path[0:3] == "std" {
		return
	}

    filepath.Ext(path)

	logDir, err := filepath.Abs(filepath.Dir(path))
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatal("SetLogout ERROR: ", err)
	}
    fileName := filepath.Join(logDir, LogName(filepath.Base(path)))
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Info("SetLogout log path : ", fileName)
		log.Out = logFile
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}

// Finalize -  log recover
func Finalize() {
	defer func() {
		err := recover()
		if err != nil {
			entry := err.(*logrus.Entry)
			log.WithFields(logrus.Fields{
				"omg":         true,
				"err_animal":  entry.Data["animal"],
				"err_size":    entry.Data["size"],
				"err_level":   entry.Level,
				"err_message": entry.Message,
				"number":      100,
			}).Error("The ice breaks!") // or use Fatal() to force the process to exit with a nonzero code
		}

	}()
	log.WithFields(logrus.Fields{
		"animal": "dcs",
	}).Info("Finish call.")
}
