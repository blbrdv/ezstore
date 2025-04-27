package msstore

import (
	"errors"
	"fmt"
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
		SetRetryWaitTime(5*time.Second).
		SetRetryMaxWaitTime(20*time.Second).
		SetHeader("Content-Encoding", "Encoding.UTF8")
}

const maxAttempts = 5

// execute [resty.Request] retrying on panic 5 times.
func execute(method string, url string, r *resty.Request) (*resty.Response, error) {
	var result *resty.Response
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt += 1 {
		func() {
			defer func() {
				if r := recover(); r != nil {
					switch recoverType := r.(type) {
					case string:
						err = errors.New(recoverType)
					case error:
						err = recoverType
					default:
						err = errors.New("unexpected type")
					}
				}
			}()

			result, err = r.Execute(method, url)
		}()

		if err == nil {
			return result, nil
		}

		pterm.Warning.Printfln("%s, Attempt %d", err.Error(), attempt)

		if attempt < maxAttempts {
			duration, _ := time.ParseDuration(fmt.Sprintf("%ds", 5*attempt))
			time.Sleep(duration)
		}
	}

	return nil, err
}
