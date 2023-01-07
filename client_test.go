package delve_client

import (
	"fmt"
	"log"
	"testing"

	"delve-client/xdlv"
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
