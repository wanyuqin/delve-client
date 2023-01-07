package xdlv

import (
	"fmt"
	"log"
	"testing"
)

func TestNewDlvBuilder(t *testing.T) {
	builder := NewDlvBuilder()
	dlv, err := builder.Exec().AddSource("./example/01/main").AddFlag("--headless").Build()
	if err != nil {
		log.Fatalf("build dlv  failed: %v", err)
	}
	err = dlv.RunServer()
	if err != nil {
		log.Fatalf("build dlv  failed: %v", err)
	}
	fmt.Println(dlv.GetServerAddr())
}
