package config

import (
	"fmt"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Initialize Logrus
func InitLogger() {
	// Convert log level from string to Logrus level
	level, err := logrus.ParseLevel(General.Log_Level)
	if err != nil {
		fmt.Println("Invalid log level in config, defaulting to INFO")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Set log format (JSON or Text)
	if General.Log_Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				// Customize how file and function name appear
				filename := fmt.Sprintf("%s:%d", f.File, f.Line)
				return "", filename
			},
			ForceColors: true,
		})
	}

	logrus.SetReportCaller(true)

	// Set output to stdout (you can also log to a file)
	logrus.SetOutput(os.Stdout)
}
