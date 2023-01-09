package delve_client

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/wanyuqin/delve-client/xdlv"
)

func init() {
}

func TestDelveClient_Connect(t *testing.T) {
	dlv, err := xdlv.NewDlvBuilder().Exec().AddSource("./example/01/main").AddFlag("--headless").Build()
	err = dlv.RunServer()
	if err != nil {
		log.Fatalf("run dlv server failed: %v", err)
	}

	client := NewClient(dlv.GetServerAddr())
	err = client.Connect()

	if err != nil {
		log.Fatalf("connect delve rpc client failed: %v", err)
	}

	_, err = client.initFuncBreakPoint()
	if err != nil {
		log.Fatalf("init func breakpoint failed: %v", err)
	}

	source, err := client.ListSource("")
	if err != nil {
		log.Fatalf("list source failed: %v", err)
	}

	for _, v := range source {
		fmt.Println(v)
	}
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
	client := NewClient("127.0.0.1:59432")
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
