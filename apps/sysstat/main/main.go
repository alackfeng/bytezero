package main

import (
	"fmt"
	"os"

	"github.com/alackfeng/bytezero/apps/sysstat"
)

func main() {
	fmt.Println(">>>>>bytezero sysstat main: pid", os.Getpid())

	stat := sysstat.NewSysStat()
	stat.Init()
	stat.Execute()
	os.Exit(1)
}
