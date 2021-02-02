// Copyright 2020 BlueCat Networks. All rights reserved

package logging

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var logger = log.New()

// AppConfig Application configuration file
type AppConfig struct {
	Logging struct {
		Level    string `yaml:"level"`
		FileName string `yaml:"file"`
	}
}

func loadAppConfig() *AppConfig {
	appConf := AppConfig{}
	yamlFile, err := ioutil.ReadFile("./app.yml")
	if err != nil {
		logger.Debugf("Failed to read the configuration file: #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &appConf)
	if err != nil {
		logger.Debugf("Failed to load the configuration: %v", err)
	}
	return &appConf
}

func getLogLevel(levelStr string) log.Level {
	var level log.Level
	switch levelStr {
	case "debug":
		level = log.DebugLevel
	case "info":
		level = log.InfoLevel
	case "warn":
		level = log.WarnLevel
	case "error":
		level = log.ErrorLevel
	default:
		level = log.WarnLevel
	}
	return level
}

// GetLogger Get the custom logger
func GetLogger() *log.Logger {
	appConf := loadAppConfig()

	var level = getLogLevel(appConf.Logging.Level)

	var fileName = "provider_bluecat.log"
	if len(appConf.Logging.FileName) > 0 {
		fileName = appConf.Logging.FileName
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	} else {
		formatter := &log.JSONFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				s := strings.Split(f.Function, ".")
				_, filename := path.Split(f.File)
				return s[len(s)-1], fmt.Sprintf("%s:%d", filename, f.Line)
			},
		}
		formatter.TimestampFormat = "02-01-2006 15:04:05"
		logger.SetReportCaller(true)
		logger.SetFormatter(formatter)
		logger.SetOutput(f)
		logger.SetLevel(level)
		return logger
	}
	return nil
}
