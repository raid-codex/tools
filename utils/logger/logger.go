package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	if os.Getenv("DEBUG") == "1" {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
