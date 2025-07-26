package windows

import (
	"bytes"
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const exeFilename = "powershell.exe"
const newline = "\r\n"

const (
	boundaryPrefix            = "$command"
	boundaryPrefixLen         = 8
	boundaryRandomPartByteLen = 12
)

type Powershell struct {
	codePage       int
	enc            encoding.Encoding
	cmd            *exec.Cmd
	stdin          io.WriteCloser
	stdout         io.ReadCloser
	stderr         io.ReadCloser
	boundaryRndBuf [boundaryRandomPartByteLen]byte
	boundaryBuf    [boundaryPrefixLen + 2*boundaryRandomPartByteLen]byte
}

var (
	ErrPowershellNotFound  = errors.New("powershell.exe not found")
	ErrUnsupportedCodePage = errors.New("unsupported code page")
)

var Encodings = map[int]encoding.Encoding{
	1200:  unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
	1201:  unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	1252:  charmap.Windows1252,
	65001: unicode.UTF8,
}

func NewPowerShell(params ...string) (*Powershell, error) {
	exePath, err := exec.LookPath(exeFilename)
	if err != nil {
		return nil, ErrPowershellNotFound
	}

	var cmd *exec.Cmd
	if len(params) > 0 {
		cmd = exec.Command(exePath, params...)
	} else {
		cmd = exec.Command(exePath, "-NoLogo", "-NoExit", "-NoProfile", "-Command", "-")
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("start powershell: %w", err)
	}

	s := &Powershell{
		enc:    encoding.Nop,
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
	copy(s.boundaryBuf[:], boundaryPrefix)

	cp, err := s.detectCodePage()
	if err != nil {
		return nil, err
	}

	enc := Encodings[cp]
	if enc == nil {
		return nil, ErrUnsupportedCodePage
	}

	s.codePage = cp
	s.enc = enc
	return s, nil
}

func (s *Powershell) CodePage() int {
	return s.codePage
}

func (s *Powershell) detectCodePage() (int, error) {
	out, err := s.Exec("[System.Text.Encoding]::Default.CodePage")
	if err != nil {
		return 0, fmt.Errorf("get codepage: %s", err.Error())
	}
	out = strings.TrimRight(out, " \r\n")
	cp, err := strconv.Atoi(out)
	if err != nil {
		return 0, fmt.Errorf("non-numeric codepage: '%s'", out)
	}
	return cp, nil
}

func (s *Powershell) Exec(cmd string) (stdout string, err error) {
	// wrap the command in special markers so we know when to stop reading from the pipes
	boundary := s.randomBoundary()
	full := fmt.Sprintf("%s; echo '%s'; [Console]::Error.WriteLine('%s')%s", cmd, boundary, boundary, newline)
	full, err = s.enc.NewEncoder().String(full)
	if err != nil {
		return "", fmt.Errorf("encode command: %s", err)
	}
	_, err = s.stdin.Write([]byte(full))
	if err != nil {
		return "", fmt.Errorf("write command: %s", err)
	}

	var stderr string
	var wg sync.WaitGroup
	wg.Add(2)
	go readOutput(s.stdout, s.enc.NewDecoder(), &stdout, boundary, &wg)
	go readOutput(s.stderr, s.enc.NewDecoder(), &stderr, boundary, &wg)
	wg.Wait()
	if len(stderr) > 0 {
		return "", errors.New(stderr)
	}
	return stdout, nil
}

func (s *Powershell) Exit() error {
	_, err := s.stdin.Write([]byte("exit" + newline))
	if err != nil {
		return fmt.Errorf("write exit: %s", err)
	}

	err = s.stdin.Close()
	if err != nil {
		return fmt.Errorf("close stdin: %s", err)
	}

	return nil
}

func (s *Powershell) randomBoundary() string {
	_, err := crand.Read(s.boundaryRndBuf[:])
	if err != nil {
		panic(err)
	}
	hex.Encode(s.boundaryBuf[boundaryPrefixLen:], s.boundaryRndBuf[:])
	return string(s.boundaryBuf[:])
}

func readOutput(r io.Reader, dec *encoding.Decoder, out *string, boundary string, wg *sync.WaitGroup) {
	var bout []byte
	defer func() {
		*out = string(bout)
		wg.Done()
	}()

	marker := []byte(boundary + newline)
	const bufsize = 64
	buf := make([]byte, bufsize)
	for {
		n, err := r.Read(buf)
		if err != nil {
			return
		}

		decoded, err := dec.Bytes(buf[:n])
		if err != nil {
			return
		}

		bout = append(bout, decoded...)
		if bytes.HasSuffix(bout, marker) {
			bout = bout[:len(bout)-len(marker)]
			return
		}
	}
}
