package xdlv

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

var apiServerPrefix = "API server listening at: "

type Dlv struct {
	cmd *exec.Cmd

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

	serverAddr string
}

type DlvBuilder struct {
	args   []string
	usage  string
	cmd    string
	flag   []string
	source string
}

var defaultUsage = "dlv"

var (
	ExecCommand  = "exec"
	DebugCommand = "debug"
)

func NewDlvBuilder() *DlvBuilder {
	return &DlvBuilder{
		usage: defaultUsage,
	}

}
func (d *DlvBuilder) Exec() *DlvBuilder {
	d.cmd = ExecCommand
	return d
}

func (d *DlvBuilder) Debug() *DlvBuilder {
	d.cmd = DebugCommand
	return d
}

func (d *DlvBuilder) AddSource(source string) *DlvBuilder {
	d.source = source
	return d
}

func (d *DlvBuilder) AddFlag(flag string) *DlvBuilder {
	d.flag = append(d.flag, flag)
	return d
}

func (d *DlvBuilder) Build() (*Dlv, error) {
	// 查找系统中是否存在盖可执行文件
	path, err := exec.LookPath(d.usage)
	if err != nil {
		return nil, err
	}
	fmt.Printf("dlv is available at %s\n", path)
	cmd := exec.Command(d.usage, d.cmd, d.source)

	if d.cmd == "" {
		return nil, errors.New("command cannot null")
	}

	if d.source == "" {
		return nil, errors.New("source cannot null")
	}
	for _, v := range d.flag {
		cmd.Args = append(cmd.Args, v)
	}
	dlv := &Dlv{
		cmd: cmd,
	}
	dlv.stdout, _ = dlv.cmd.StdoutPipe()
	dlv.stdin, _ = dlv.cmd.StdinPipe()
	dlv.stderr, _ = dlv.cmd.StderrPipe()
	return dlv, nil
}

func (d *Dlv) RunServer() error {
	err := d.cmd.Start()

	if err != nil {
		return err
	}
	d.parseListeningAddr()
	if d.serverAddr != "" {
		fmt.Printf("dlv server run at %s \n", d.serverAddr)
	}
	return err
}

func (d *Dlv) GetServerAddr() string {
	return d.serverAddr
}

func (d *Dlv) parseListeningAddr() {
	buf := make([]byte, 4*1024)
	flag := true

	for {
		n, _ := d.stdout.Read(buf)

		text := string(buf[:n])
		if flag && strings.HasPrefix(text, apiServerPrefix) {
			nl := strings.Index(text, "\n")
			if nl > 0 {
				text = text[:nl]
				d.serverAddr = text[len(apiServerPrefix):]
				flag = false
				return
			}

		}
	}
}
