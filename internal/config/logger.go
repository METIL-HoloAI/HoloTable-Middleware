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

	//customize log output
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// Currently, we print the filepath and line number
			filename := fmt.Sprintf("%s:%d", f.File, f.Line)
			return "", filename
		},
		ForceColors: true,
	})

	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)
}
