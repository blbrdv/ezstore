package msstore

import (
	"errors"
	"github.com/pterm/pterm"
	"time"

	"github.com/go-resty/resty/v2"
)

type ptermLogger struct {
	resty.Logger
}

func (l ptermLogger) Errorf(format string, v ...interface{}) {
	pterm.Error.Printfln(format, v...)
}

func (l ptermLogger) Warnf(format string, v ...interface{}) {
	pterm.Warning.Printfln(format, v...)
}

func (l ptermLogger) Debugf(format string, v ...interface{}) {
	pterm.Debug.Printfln(format, v...)
}

func http() *resty.Client {
	return resty.
		New().
		SetLogger(ptermLogger{}).
		SetRetryCount(10).
		SetRetryWaitTime(5*time.Second).
		SetRetryMaxWaitTime(20*time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		}).
		SetHeader("Content-Encoding", "Encoding.UTF8")
}
