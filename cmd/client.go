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
	"fmt"

	"github.com/alackfeng/bytezero/apps/client"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client called")
        tcpAddress, _ := cmd.Flags().GetString("tcp-address")
        maxBufferLen, _ := cmd.Flags().GetInt("max-buffer-len")
        sendPeroidMs, _ := cmd.Flags().GetInt("send-peroid-ms")
        recvCheck, _ := cmd.Flags().GetBool("recv-check")


        appsClient := client.NewAppsClient(tcpAddress)
        appsClient.SetMaxBufferLen(maxBufferLen).SetSendPeroid(sendPeroidMs).SetRecvCheck(recvCheck).Main()
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringP("tcp-address", "t", "127.0.0.1:7788", "TCP Address")
	clientCmd.Flags().IntP("max-buffer-len", "l", 1024*1024*1, "Max Buffer Length")
	clientCmd.Flags().IntP("send-peroid-ms", "p", 10, "Send Peroid Ms.")
	clientCmd.Flags().BoolP("recv-check", "c", false, "Recv Check, false is mean to close connection.")
}
