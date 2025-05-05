package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/pterm/pterm"
	"io"
	"os"
	"path/filepath"
	"strings"

	"time"

	"github.com/imroc/req/v3"
)

func getDumpFile() *os.File {
	cache, _ := os.UserCacheDir()
	filename := filepath.Join(cache, "ezstore", *log.GetLogFileName())
	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)

	return file
}

var dumpFileName = getDumpFile()

func traceRequest(req *req.Request) {
	if req == nil {
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(dumpFileName)

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

	_, _ = dumpFileName.WriteString(sb.String())
}

func traceError(err error) {
	if err == nil {
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(dumpFileName)

	var sb strings.Builder

	text := err.Error()

	if text == "" {
		_, _ = fmt.Fprintln(&sb, "! (empty)")
	} else {
		_, _ = fmt.Fprintf(&sb, "! %s\n", text)
	}

	_, _ = fmt.Fprint(&sb, "\n")

	_, _ = dumpFileName.WriteString(sb.String())
}

func traceResponse(res *req.Response) {
	if res == nil || res.Response == nil {
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(dumpFileName)

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

	_, _ = dumpFileName.WriteString(sb.String())
}

func getHTTPClient() *req.Client {
	client := req.C().
		SetCommonHeader("Content-Encoding", "Encoding.UTF8").
		SetCommonRetryCount(5).
		SetCommonRetryInterval(func(_ *req.Response, attempt int) time.Duration {
			result, _ := time.ParseDuration(fmt.Sprintf("%ds", 5*attempt))
			return result
		}).
		AddCommonRetryHook(func(_ *req.Response, err error) {
			if err != nil {
				pterm.Warning.Println(err.Error())
			}
		})

	if log.IsTraceLevel() {
		client = client.
			OnError(func(_ *req.Client, req *req.Request, resp *req.Response, err error) {
				traceRequest(req)
				traceError(err)
				traceResponse(resp)
			}).
			OnAfterResponse(func(_ *req.Client, resp *req.Response) error {
				traceRequest(resp.Request)
				traceError(resp.Err)
				traceResponse(resp)
				return nil
			})
	}

	return client
}

var client = getHTTPClient()
