package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/blbrdv/ezstore/internal/ms/windows"
	"github.com/blbrdv/ezstore/internal/utils"
	"io"
	"os"
	"strings"

	"time"

	"github.com/imroc/req/v3"
)

func getTraceFile() *windows.File {
	return windows.OpenFile(log.TraceFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE)
}

func traceRequest(traceFile *windows.File, req *req.Request) {
	if req == nil {
		return
	}

	var sb strings.Builder

	utils.Fprintf(&sb, "> %s %s%s", strings.ToUpper(req.Method), req.RawURL, utils.WindowsNewline)

	for key, values := range req.Headers {
		for _, value := range values {
			utils.Fprintf(&sb, "> %s: %s%s", key, value, utils.WindowsNewline)
		}
	}

	if req.Body != nil {
		utils.Fprintln(&sb, "> Body:")
		utils.Fprintf(&sb, "%s%s", req.Body, utils.WindowsNewline)
	}

	utils.Fprint(&sb, utils.WindowsNewline)

	traceFile.WriteString(sb.String())
}

func traceError(traceFile *windows.File, err error) {
	if err == nil {
		return
	}

	var sb strings.Builder

	text := err.Error()

	if text == "" {
		utils.Fprintln(&sb, "! (empty)")
	} else {
		utils.Fprintf(&sb, "! %s%s", text, utils.WindowsNewline)
	}

	utils.Fprint(&sb, utils.WindowsNewline)

	traceFile.WriteString(sb.String())
}

func traceResponse(traceFile *windows.File, res *req.Response) {
	if res == nil || res.Response == nil {
		return
	}

	var sb strings.Builder

	status := res.GetStatus()
	if status != "" {
		utils.Fprintf(&sb, "< %s%s", status, utils.WindowsNewline)
	}

	if len(res.Header) > 0 {
		for key, values := range res.Header {
			for _, value := range values {
				utils.Fprintf(&sb, "< %s: %s%s", key, value, utils.WindowsNewline)
			}
		}
	}

	bodyRaw, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}
	if len(bodyRaw) > 0 {
		utils.Fprintln(&sb, "< Body:")

		if body := string(bodyRaw); body == "" {
			for _, data := range bodyRaw {
				utils.Fprintf(&sb, "%b", data)
			}
		} else {
			utils.Fprintf(&sb, "%s%s", body, utils.WindowsNewline)
		}
	}

	utils.Fprint(&sb, utils.WindowsNewline)

	traceFile.WriteString(sb.String())
}

func either(e1 error, e2 error) error {
	if e1 != nil {
		return e1
	} else if e2 != nil {
		return e2
	} else {
		return nil
	}
}

func getHTTPClient() *req.Client {
	return req.C().
		SetCommonHeader("Content-Encoding", "Encoding.UTF8").
		SetCommonRetryCount(5).
		SetCommonRetryInterval(func(_ *req.Response, attempt int) time.Duration {
			result, err := time.ParseDuration(fmt.Sprintf("%ds", 5*attempt))
			if err != nil {
				panic(err.Error())
			}
			return result
		}).
		SetCommonRetryHook(func(resp *req.Response, err error) {
			underlyingErr := either(err, resp.Err)
			if log.Level == log.Detailed {
				file := getTraceFile()
				traceRequest(file, resp.Request)
				traceError(file, underlyingErr)
				traceResponse(file, resp)
				file.Close()
			}

			if underlyingErr == nil {
				log.Warningf("%s %s: %s", strings.ToUpper(resp.Request.Method), resp.Request.RawURL, resp.GetStatus())
			} else {
				log.Warning(underlyingErr.Error())
			}
		}).
		OnError(func(_ *req.Client, req *req.Request, resp *req.Response, err error) {
			if log.Level == log.Detailed {
				underlyingErr := either(err, resp.Err)
				file := getTraceFile()
				traceRequest(file, req)
				traceError(file, underlyingErr)
				traceResponse(file, resp)
				file.Close()
			}
		})
}

var client = getHTTPClient()
