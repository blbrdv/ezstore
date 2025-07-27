package windows

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
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

type cmd struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func (c *cmd) exec(cmd string, nl string, enc encoding.Encoding) (stdout string, err error) {
	// wrap the command in special markers so we know when to stop reading from the pipes
	boundary := "debugdebug"
	full := fmt.Sprintf("%s; echo '%s'; [Console]::Error.WriteLine('%s')%s", cmd, boundary, boundary, nl)
	log.Tracef("Full cmd: %s", full) // TODO: remove after tests
	full, err = enc.NewEncoder().String(full)
	if err != nil {
		return "", fmt.Errorf("encode command: %s", err)
	}
	log.Trace("Command encoded") // TODO: remove after tests
	_, err = c.stdin.Write([]byte(full))
	if err != nil {
		return "", fmt.Errorf("write command: %s", err)
	}
	log.Trace("Stdin filled") // TODO: remove after tests

	var stderr string
	var wg sync.WaitGroup
	wg.Add(2)
	go readOutput("out", c.stdout, enc.NewDecoder(), &stdout, boundary, nl, &wg)
	go readOutput("err", c.stderr, enc.NewDecoder(), &stderr, boundary, nl, &wg)
	wg.Wait()
	log.Trace("Both readOutput finished") // TODO: remove after tests
	if len(stderr) > 0 {
		return "", errors.New(stderr)
	}
	return stdout, nil
}

func (c *cmd) Start() error {
	return c.cmd.Start()
}

type Powershell struct {
	*cmd
	codePage int
	enc      encoding.Encoding
	newLine  string
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

var (
	defaultParams = []string{
		"-NoLogo",
		"-NoProfile",
	}
	noExitParam = "-NoExit"
	cmdParams   = []string{
		"-Command",
		"-",
	}
)

func NewPowerShell(params ...string) (*Powershell, error) {
	exePath, err := exec.LookPath(exeFilename)
	if err != nil {
		return nil, ErrPowershellNotFound
	}
	log.Tracef("Found exe: '%s'", exePath) // TODO: remove after tests

	var settingsParams []string
	var commandParams []string
	if len(params) > 0 {
		settingsParams = append(params, cmdParams...)
		commandParams = append(params, noExitParam)
		commandParams = append(commandParams, cmdParams...)
	} else {
		settingsParams = append(defaultParams, cmdParams...)
		commandParams = append(defaultParams, noExitParam)
		commandParams = append(commandParams, cmdParams...)
	}
	log.Trace("Params set") // TODO: remove after tests

	nl, cp, err := getDefaultSettings(exePath, settingsParams...)
	if err != nil {
		return nil, err
	}

	execCmd, err := getCmd(exePath, commandParams...)
	if err != nil {
		return nil, err
	}

	err = execCmd.Start()
	if err != nil {
		return nil, fmt.Errorf("start powershell: %w", err)
	}
	log.Trace("Started PowerShell") // TODO: remove after tests

	enc := Encodings[cp]
	if enc == nil {
		return nil, ErrUnsupportedCodePage
	}
	log.Trace("Encoding set") // TODO: remove after tests

	return &Powershell{
		cmd:      execCmd,
		codePage: cp,
		enc:      encoding.Nop,
		newLine:  nl,
	}, nil
}

func (s *Powershell) CodePage() int {
	return s.codePage
}

func (s *Powershell) Exec(cmd string) (string, error) {
	return s.exec(cmd, s.newLine, s.enc)
}

func (s *Powershell) Exit() error {
	_, err := s.stdin.Write([]byte("exit" + s.newLine))
	if err != nil {
		return fmt.Errorf("write exit: %s", err)
	}

	err = s.stdin.Close()
	if err != nil {
		return fmt.Errorf("close stdin: %s", err)
	}

	return nil
}

func readOutput(name string, r io.Reader, dec *encoding.Decoder, out *string, boundary string, nl string, wg *sync.WaitGroup) {
	var bout []byte
	defer func() {
		*out = string(bout)
		wg.Done()
	}()

	marker := []byte(boundary + nl)
	const bufsize = 64
	buf := make([]byte, bufsize)
	log.Tracef("%s readOutput started", name) // TODO: remove after tests
	for {
		n, err := r.Read(buf)
		if err != nil {
			return
		}
		log.Tracef("%s: read buff: %d", name, n) // TODO: remove after tests

		decoded, err := dec.Bytes(buf[:n])
		if err != nil {
			return
		}
		log.Tracef("%s: decoded: %+q", name, decoded) // TODO: remove after tests

		bout = append(bout, decoded...)
		if bytes.HasSuffix(bout, marker) {
			bout = bout[:len(bout)-len(marker)]
			return
		}
		log.Tracef("%s: bout: %+q", name, bout) // TODO: remove after tests
	}
}

func getDefaultSettings(exePath string, params ...string) (string, int, error) {
	codepageParams := append(params, "[System.Text.Encoding]::Default.CodePage")
	out, err := exec.Command(exePath, codepageParams...).Output()
	if err != nil {
		return "", 0, err
	}
	codepage := string(out)
	log.Tracef("Codepage recieved: %s", codepage) // TODO: remove after tests
	codepage = strings.TrimRight(codepage, " "+newline)
	cp, err := strconv.Atoi(codepage)
	if err != nil {
		return "", 0, fmt.Errorf("non-numeric codepage: '%s'", codepage)
	}
	log.Tracef("Codepage converted: %d", cp) // TODO: remove after tests

	nlParams := append(params, "[Environment]::NewLine -replace '`',\"\\\"")
	out, err = exec.Command(exePath, nlParams...).Output()
	if err != nil {
		return "", 0, err
	}
	nl := string(out)
	log.Tracef("NewLine recieved: %s", nl) // TODO: remove after tests

	return nl, cp, nil
}

func getCmd(exePath string, params ...string) (*cmd, error) {
	execCmd := exec.Command(exePath, params...)
	log.Trace("Set cmd") // TODO: remove after tests

	stdin, err := execCmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	log.Trace("Set StdinPipe") // TODO: remove after tests

	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	log.Trace("Set StdoutPipe") // TODO: remove after tests

	stderr, err := execCmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}
	log.Trace("Set StderrPipe") // TODO: remove after tests

	return &cmd{cmd: execCmd, stdin: stdin, stdout: stdout, stderr: stderr}, nil
}
