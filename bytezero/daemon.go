package bytezero

import (
	"fmt"
	"os"
	"os/exec"
)

// IsDaemon -
func IsDaemon() bool {
    args := os.Args
    daemon := false
    for k, v := range args {
        if v == "--daemon" {
	    daemon = true
	    args[k] = ""
	}
    }
    return daemon
}

// Daemon -
func Daemon() error {
    cmdName := os.Args[0]
    cmdArgs := os.Args[1:]
    fmt.Println(cmdName, cmdArgs, " .")
    cmd := exec.Command(cmdName, cmdArgs...)
    err := cmd.Start()
    if err != nil {
        return err
    }
    fmt.Println("bytezero Daemon mode, process id ", cmd.Process.Pid)
    os.Exit(0)
    return nil
}
