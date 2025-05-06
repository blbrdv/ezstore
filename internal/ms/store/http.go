package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"io"
	"os"
	"strings"

	"time"

	"github.com/imroc/req/v3"
)

func getTraceFile() *os.File {
	file, err := os.OpenFile(log.TraceFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Sprintf("can not open trace file %s: %s", log.TraceFile, err.Error()))
	}

	return file
}

func traceRequest(traceFile *os.File, req *req.Request) {
	if req == nil {
		return
	}

	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "> %s %s\n", strings.ToUpper(req.Method), req.RawURL)

	for key, values := range req.Headers {
		for _, value := range values {
			_, _ = fmt.Fprintf(&sb, "> %s: %s\n", key, value)
		}
	}

	if req.Body != nil {
		_, _ = fmt.Fprintln(&sb, "> Body:")
		_, _ = fmt.Fprintf(&sb, "%s\n", req.Body)
	}

	_, _ = fmt.Fprint(&sb, "\n")

	_, _ = traceFile.WriteString(sb.String())
}

func traceError(traceFile *os.File, err error) {
	if err == nil {
		return
	}

	var sb strings.Builder

	text := err.Error()

	if text == "" {
		_, _ = fmt.Fprintln(&sb, "! (empty)")
	} else {
		_, _ = fmt.Fprintf(&sb, "! %s\n", text)
	}

	_, _ = fmt.Fprint(&sb, "\n")

	_, _ = traceFile.WriteString(sb.String())
}

func traceResponse(traceFile *os.File, res *req.Response) {
	if res == nil || res.Response == nil {
		return
	}

	var sb strings.Builder

	status := res.Status
	if status != "" {
		_, _ = fmt.Fprintf(&sb, "< %s\n", status)
	}

	if len(res.Header) > 0 {
		for key, values := range res.Header {
			for _, value := range values {
				_, _ = fmt.Fprintf(&sb, "< %s: %s\n", key, value)
			}
		}
	}

	bodyRaw, _ := io.ReadAll(res.Body)
	if len(bodyRaw) > 0 {
		_, _ = fmt.Fprintln(&sb, "< Body:")

		if body := string(bodyRaw); body == "" {
			for _, data := range bodyRaw {
				_, _ = fmt.Fprintf(&sb, "%b", data)
			}
		} else {
			_, _ = fmt.Fprintf(&sb, "%s\n", body)
		}
	}

	_, _ = fmt.Fprint(&sb, "\n")

	_, _ = traceFile.WriteString(sb.String())
}

func getHTTPClient() *req.Client {
	return req.C().
		SetCommonHeader("Content-Encoding", "Encoding.UTF8").
		SetCommonRetryCount(5).
		SetCommonRetryInterval(func(_ *req.Response, attempt int) time.Duration {
			result, _ := time.ParseDuration(fmt.Sprintf("%ds", 5*attempt))
			return result
		}).
		AddCommonRetryHook(func(resp *req.Response, err error) {
			if err != nil {
				if log.Level == log.Detailed {
					file := getTraceFile()
					traceRequest(file, resp.Request)
					traceError(file, resp.Err)
					traceResponse(file, resp)
					_ = file.Close()
				}

				log.Warning(err.Error())
			}
		}).
		OnError(func(_ *req.Client, req *req.Request, resp *req.Response, err error) {
			if log.Level == log.Detailed {
				file := getTraceFile()
				traceRequest(file, req)
				traceError(file, err)
				traceResponse(file, resp)
				_ = file.Close()
			}
		})
}

var client = getHTTPClient()
