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
    for _, v := range args {
        if v == "--daemon" {
            daemon = true
            break // args[k] = ""
        }
    }
    return daemon
}

const BYTEZERO_DAEMON = "BYTEZERO_DAEMON"
// Daemon -
func Daemon() error {
    if os.Getenv(BYTEZERO_DAEMON) == "1" {
        return nil
    }
    path := os.Args[0]
    args := os.Args[1:]
    fmt.Println(path, args, " .")
    cmd := exec.Command(path, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", BYTEZERO_DAEMON, "1"))
    err := cmd.Start()
    if err != nil {
        return err
    }
    fmt.Println("bytezero Daemon mode, process id ", cmd.Process.Pid)
    os.Exit(0) // quit parent process.
    return nil
}
