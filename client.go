package delve_client

import (
	"errors"

	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/rpc2"
)

type DelveClient struct {
	Addr      string
	C         *rpc2.RPCClient
	StartFunc string
}

var ClientIsNil = errors.New("RPCClient is nil")

var startFunc = "main.main"

func NewClient(addr string) *DelveClient {

	return &DelveClient{
		Addr: addr,
	}
}

func (d *DelveClient) Connect() error {
	c := rpc2.NewClient(d.Addr)
	if c == nil {
		return ClientIsNil
	}
	d.C = c

	return nil
}

func (d *DelveClient) Halt() (*api.DebuggerState, error) {
	return d.C.Halt()
}

// ListSource 获取所有源代码
func (d *DelveClient) ListSource(filter string) ([]string, error) {
	return d.C.ListSources(filter)
}

func (d *DelveClient) ListFunctions(filter string) ([]string, error) {
	return d.C.ListFunctions(filter)
}

func (d *DelveClient) initFuncBreakPoint() (int, error) {
	if d.StartFunc == "" {
		d.StartFunc = startFunc
	}

	b := api.Breakpoint{FunctionName: startFunc}

	r, err := d.C.CreateBreakpoint(&b)

	return r.ID, err
}
