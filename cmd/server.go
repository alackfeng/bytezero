/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bze "github.com/alackfeng/bytezero/cores"
	profile "github.com/alackfeng/bytezero/cores/utils"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bytezero server called", daemonProc, bze.ConfigGlobal().App.LogPath)

        profile.SetLogout(bze.ConfigGlobal().App.LogPath)
        profile.InitGC(100)
        stopFunc, err := profile.ProfileIfEnabled()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bytezero Server Error: %s\n", err.Error())
			return
		}
		defer stopFunc() // to be executed as late as possible


        maxBufferLen, _ := cmd.Flags().GetInt("max-buffer-len")
        rwBufferLen, _ := cmd.Flags().GetInt("rw-buffer-len")
        port, _ := cmd.Flags().GetInt("port")
        margic, _ := cmd.Flags().GetBool("margic")
        host, _ := cmd.Flags().GetString("host")
        appid, _ := cmd.Flags().GetString("appid")
        appkey, _ := cmd.Flags().GetString("appkey")
        needTls, _ := cmd.Flags().GetBool("tls")
        tlsPort, _ := cmd.Flags().GetInt("tlsport")
        caCert, _ := cmd.Flags().GetString("cacert")
        caKey, _ := cmd.Flags().GetString("cakey")

        done := make(chan bool)
		sigs := make(chan os.Signal, 1)
        signal.Notify(sigs)
        ctx, cancel := context.WithCancel(context.Background())
        defer func() {
			cancel()
			close(sigs)
			close(done)
		}()
        bzn := bze.NewBytezeroNet(ctx, done)
        bze.ConfigSetServer(maxBufferLen, rwBufferLen, port, host, appid, appkey, margic)
        bze.ConfigSetTls(needTls, tlsPort, caCert, caKey)
        bzn.Main()

        logcmd.Errorln("main listen...")
		bQuit := false
		for {
			select {
			case sig := <-sigs:
				if sig == syscall.SIGTERM || sig == syscall.SIGINT {
					logcmd.Warnln("main Catch Signal: ", sig)
					bQuit = bzn.Quit()
					// } else if sig == syscall.SIGURG {
				} else {
					if sig.String() == "child exited" || sig.String() == "urgent I/O condition" {
					} else {
						logcmd.Warnln("main Catch Signal - ", sig)
					}
				}
			case d := <-done:
				logcmd.Errorln("main done. ", d)
				bQuit = true
			}

			if bQuit { // QUIT
				break
			}
		}
		logcmd.Errorln("main over...")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IntP("max-buffer-len", "l", 1024*1024*10, "Max Buffer Length")
	serverCmd.Flags().IntP("rw-buffer-len", "b", 1024*1024*1, "Read and Write Buffer Length")
	serverCmd.Flags().IntP("port", "p", 7788, "tcp or udp server listen port")
    serverCmd.Flags().BoolP("margic", "m", true, "margic for tcp listen secret")
	serverCmd.Flags().StringP("host", "s", "192.168.90.162:7790", "web rest api host url")
    serverCmd.Flags().StringP("appid", "i", "bytezero-appid", "bytezero appid")
	serverCmd.Flags().StringP("appkey", "a", "secret", "appkey secret")
    serverCmd.Flags().BoolP("tls", "t", true, "tls server listen")
    serverCmd.Flags().IntP("tlsport", "r", 7789, "tls server listen port")
    serverCmd.Flags().StringP("cacert", "e", "./scripts/certs/server/server.crt", "tls server ca cert: ./scripts/certs/server/server.crt")
    serverCmd.Flags().StringP("cakey", "k", "./scripts/certs/server/server.key", "tls server ca key: ./scripts/certs/server/server.key")
}
