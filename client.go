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

// NewClient 创建一个dlv客户端
func NewClient(addr string) *DelveClient {
	return &DelveClient{
		Addr: addr,
	}
}

// Connect 连接到Dlv服务
func (d *DelveClient) Connect() error {
	c := rpc2.NewClient(d.Addr)
	if c == nil {
		return ClientIsNil
	}
	d.C = c

	return d.initMainFuncBreakPoint()
}

func (d *DelveClient) Halt() (*api.DebuggerState, error) {
	return d.C.Halt()
}

// ListSource 获取所有源代码
func (d *DelveClient) ListSource(filter string) ([]string, error) {
	return d.C.ListSources(filter)
}

// ListFunctions 列出所有的Function
func (d *DelveClient) ListFunctions(filter string) ([]string, error) {
	return d.C.ListFunctions(filter)
}

// ListBreakpoints 列出所有的断点
func (d *DelveClient) ListBreakpoints(all bool) ([]*api.Breakpoint, error) {
	return d.C.ListBreakpoints(all)
}

// ListGoroutines 列出所有的goroutines
func (d *DelveClient) ListGoroutines(start int, count int) ([]*api.Goroutine, int, error) {
	return d.C.ListGoroutines(start, count)
}

func (d *DelveClient) GetState() (*api.DebuggerState, error) {
	return d.C.GetState()
}

// CreateBreakPoints 创建断点
func (d *DelveClient) CreateBreakPoints(breakPoint *api.Breakpoint) (*api.Breakpoint, error) {
	return d.C.CreateBreakpoint(breakPoint)
}

// ClearBreakpoint 清除断点
func (d *DelveClient) ClearBreakpoint(id int) (*api.Breakpoint, error) {
	return d.C.ClearBreakpoint(id)
}

func (d *DelveClient) Stacktrace(goroutineId int64, depth int, opts api.StacktraceOptions, cfg *api.LoadConfig) ([]api.Stackframe, error) {
	return d.C.Stacktrace(goroutineId, depth, opts, cfg)
}

// Continue 运行到下一个断点结束 或者程序结束
func (d *DelveClient) Continue() <-chan *api.DebuggerState {
	return d.C.Continue()
}

func (d *DelveClient) Step() (*api.DebuggerState, error) {
	return d.C.Step()
}

// Next 单步执行
func (d *DelveClient) Next() (*api.DebuggerState, error) {
	return d.C.Next()
}

// ListLocalVariables 列出作用域中所有的本地变量
func (d *DelveClient) ListLocalVariables(scope api.EvalScope, cfg api.LoadConfig) ([]api.Variable, error) {
	return d.C.ListLocalVariables(scope, cfg)
}

func (d *DelveClient) EvalVariable(scope api.EvalScope, expr string, cfg api.LoadConfig) (*api.Variable, error) {
	return d.C.EvalVariable(scope, expr, cfg)
}

// initMainFuncBreakPoint 初始化main.main断点
// 主要是程序启动时候添加一个默认断点，并且将程序停止到main
func (d *DelveClient) initMainFuncBreakPoint() error {
	if d.StartFunc == "" {
		d.StartFunc = startFunc
	}

	b := &api.Breakpoint{FunctionName: startFunc, Line: -1}

	r, err := d.CreateBreakPoints(b)

	for range d.Continue() {

	}

	defer d.ClearBreakpoint(r.ID)

	return err
}
