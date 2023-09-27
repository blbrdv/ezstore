package msstore

import (
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
)

func http() *resty.Client {
	return resty.
		New().
		SetRetryCount(10).
		SetRetryWaitTime(5*time.Second).
		SetRetryMaxWaitTime(20*time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		}).
		SetHeader("Content-Encoding", "Encoding.UTF8")
}
