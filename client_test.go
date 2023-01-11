package delve_client

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/go-delve/delve/service/api"

	"github.com/wanyuqin/delve-client/pkg/logger"
	"github.com/wanyuqin/delve-client/xdlv"
)

var normalLoadConfig = api.LoadConfig{
	FollowPointers:     false,
	MaxVariableRecurse: 1,
	MaxStringLen:       64,
	MaxArrayValues:     64,
	MaxStructFields:    -1,
}

func TestDelveClient_Connect(t *testing.T) {
	dlv, err := xdlv.NewDlvBuilder().Exec().AddSource("./example/01/main").AddFlag("--headless").Build()
	err = dlv.RunServer()
	if err != nil {
		log.Fatalf("run dlv server failed: %v", err)
	}

	// 监听 stdout 输出
	go func() {
		for {
			select {
			case m := <-*dlv.StdoutChan:
				logger.Debugf("stdout: %s", m)
			}
		}
	}()

	go func() {
		for {
			select {
			case m := <-*dlv.StderrChan:
				logger.Debugf("stderr: %s", m)
			}
		}
	}()

	client := NewClient(dlv.GetServerAddr())
	err = client.Connect()

	if err != nil {
		log.Fatalf("connect delve rpc client failed: %v", err)
	}

	// breakpoints, err := client.ListBreakpoints(true)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, v := range breakpoints {
	// 	fmt.Println(v)
	// }

	// states = <-client.Continue()
	// states = <-client.Continue()

	// goroutines, i, err := client.ListGoroutines(0, 10)
	// if err != nil {
	// 	log.Fatalf("list goroutines failed: %v", err)
	// }
	//
	// fmt.Printf("goroutines count %d \n", i)
	// for _, v := range goroutines {
	// 	fmt.Println(v)
	// }

	for {
		_, err = client.Next()
		if err != nil {
			break
		}
		state, err := client.C.GetState()

		if err != nil {
			logger.Errorf("GetState failed: %v", err)
			break
		}

		if state.Exited {
			logger.Error("debugger state exited")
			break
		}

		// fmt.Printf("debugger state %v \n", state)
		// es := api.EvalScope{
		// 	GoroutineID: state.CurrentThread.GoroutineID,
		// }

		// variables, err := client.ListLocalVariables(es, normalLoadConfig)
		// if err != nil {
		// 	logger.Errorf("ListLocalVariables failed: %v", err)
		// }

		// for _, v := range variables {
		// 	variable, err := client.EvalVariable(es, v.Name, normalLoadConfig)
		// 	if err != nil {
		// 		logger.Errorf("EvalVariable failed: %v", err)
		// 		continue
		// 	}
		//
		// 	logger.Debugf("var %s val is %s type is %s", variable.Name, variable.Value, variable.Type)
		// }
	}

	// source, err := client.ListSource("")
	// if err != nil {
	// 	log.Fatalf("list source failed: %v", err)
	// }
	//
	// for _, v := range source {
	// 	fmt.Println(v)
	// }

	//
	// functions, err := client.ListFunctions("main.main")
	// if err != nil {
	// 	log.Fatalf("list functions failed: %v", err)
	// }
	//
	// for _, v := range functions {
	// 	fmt.Println(v)
	// }
}

func TestNewClient(t *testing.T) {
	client := NewClient("127.0.0.1:56719")
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/sources", func(w http.ResponseWriter, r *http.Request) {
		source, err := client.ListSource("")
		if err != nil {
			fmt.Fprintf(w, "sources failed: %v", err)
			return
		}

		fmt.Fprintln(w, source)
	})

	mux.HandleFunc("/source-main", func(w http.ResponseWriter, r *http.Request) {
		source, err := client.ListSource("")
		if err != nil {
			fmt.Fprintf(w, "sources failed: %v", err)
			return
		}
		if len(source) == 0 {
			fmt.Fprint(w, "cannot find main.go")
			return
		}

		sm := source[0]
		_, err = os.Stat(sm)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		fb, err := os.ReadFile(sm)

		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		fmt.Fprint(w, string(fb))
	})

	mux.HandleFunc("/breakpoints", func(w http.ResponseWriter, r *http.Request) {
		breakpoints, err := client.ListBreakpoints(true)
		if err != nil {
			fmt.Fprintf(w, "list break points failed: %v", err)
			return
		}
		for _, v := range breakpoints {
			fmt.Fprintln(w, v.FunctionName)
		}

	})

	http.ListenAndServe("127.0.0.1:8080", mux)
}
