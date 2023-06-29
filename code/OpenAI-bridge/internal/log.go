package internal

import (
	"Learning_Record/code/OpenAI-bridge/config"
	"github.com/sirupsen/logrus"
)

type Params struct {
	Name string
}

func PrintLogRus(level, funcName string, args ...interface{}) {
	switch level {
	case "trace":
		config.Log.WithFields(logrus.Fields{
			"funcName": funcName,
		}).Trace(args)
	case "debug":
		config.Log.WithFields(logrus.Fields{
			"funcName": funcName,
		}).Debug(args)
	case "info":
		config.Log.WithFields(logrus.Fields{
			"funcName": funcName,
		}).Info(args)
	case "warn":
		config.Log.WithFields(logrus.Fields{
			"funcName": funcName,
		}).Warn(args)
	case "error":
		config.Log.WithFields(logrus.Fields{
			"funcName": funcName,
		}).Error(args)
	case "fatal":
		config.Log.WithFields(logrus.Fields{
			"funcName": funcName,
		}).Fatal(args)
	case "panic":
		config.Log.WithFields(logrus.Fields{
			"funcName": funcName,
		}).Panic(args)
	default:
		config.Log.WithFields(logrus.Fields{
			"funcName": funcName,
		}).Debug(args)
	}
}
