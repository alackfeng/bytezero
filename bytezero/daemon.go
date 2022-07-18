package bytezero

import (
	"fmt"
	"os"
	"os/exec"
)

// Daemon -
func Daemon() error {
    fmt.Println("bytezero Daemon mode.")
    cmdName := os.Args[0]
    cmdArgs := os.Args[1:]
    cmd := exec.Command(cmdName, cmdArgs...)
    err := cmd.Start()
    if err != nil {
        return err
    }
    return nil
}
