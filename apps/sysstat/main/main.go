package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alackfeng/bytezero/apps/sysstat"
)

const (
	ver = "1.0.0"
)

func main() {
	fmt.Println(">>>>>bytezero sysstat main: pid", os.Getpid())

	stat := sysstat.NewSysStat()
	stat.Init()
	cancel := stat.Execute()

	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		fmt.Printf(">>>>>bytezero sysstat get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if cancel != nil {
				cancel()
			}
			fmt.Printf(">>>>>bytezero sysstat [version: %s] exit", ver)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
