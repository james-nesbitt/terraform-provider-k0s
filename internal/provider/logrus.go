package provider

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	logrus "github.com/sirupsen/logrus"
	logrus_writer "github.com/sirupsen/logrus/hooks/writer"
)

func PipeLogRusToTFLog(ctx context.Context) error {
	logrus.SetOutput(ioutil.Discard)

	for _, level := range logrus.AllLevels {
		logrus.AddHook(&logrus_writer.Hook{
			Writer:    LogRusToTFLogHandler{loglevel: level, context: ctx},
			LogLevels: []log.Level{level},
		})
	}

	return nil
}

// LogRusToTFLogHandler logrus to tflog writer handler
type LogRusToTFLogHandler struct {
	loglevel logrus.Level
	context context.Context // Yes this is considered bad, but this exists
	                        // in a limited scope and so is managed.
 }

// Write copy logrus entries to tflog using the log level
func (lh LogRusToTFLogHandler) Write(p []byte) (n int, err error) {
	switch lh.loglevel {
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		tflog.Error(lh.context, string(p), map[string]interface{})
		return len(p), nil
	case logrus.WarnLevel:
		tflog.Warn(lh.context, string(p), map[string]interface{})
		return len(p), nil
	case logrus.InfoLevel:
		tflog.Info(lh.context, string(p), map[string]interface{})
		return len(p), nil
	case logrus.DebugLevel:
		tflog.Debug(lh.context, string(p), map[string]interface{})
		return len(p), nil
	default:
		return 0, fmt.Errorf("unknown Log level")
	}
}
