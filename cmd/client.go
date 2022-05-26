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
        udpAddress, _ := cmd.Flags().GetString("udp-address")
        maxBufferLen, _ := cmd.Flags().GetInt("max-buffer-len")
        sendPeroidMs, _ := cmd.Flags().GetInt("send-peroid-ms")
        recvCheck, _ := cmd.Flags().GetBool("recv-check")
        sessionId, _ := cmd.Flags().GetString("session-id")
        deviceId, _ := cmd.Flags().GetString("device-id")

        // if err := checkMust(deviceId, tcpAddress, sessionId); err != nil {
        //     fmt.Println(err.Error())
        //     return
        // }

        appsClient := client.NewAppsClient()
        appsClient.SetTcpAddress(tcpAddress).SetUdpAddress(udpAddress)
        appsClient.SetDeviceId(deviceId).SetSessionId(sessionId)
        appsClient.SetMaxBufferLen(maxBufferLen).SetSendPeroid(sendPeroidMs).SetRecvCheck(recvCheck).Main()
	},
}

// checkMust -
func checkMust(deviceId, tcpAddress, sessionId string) error {
    if deviceId == "" {
        return fmt.Errorf("Set device-id Pararm")
    }
    if tcpAddress == "" {
        return fmt.Errorf("Set tcp-address Pararm")
    }
    if sessionId == "" {
        return fmt.Errorf("Set session-id Pararm")
    }
    return nil
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringP("tcp-address", "t", "", "TCP Address 127.0.0.1:7788")
	clientCmd.Flags().StringP("udp-address", "u", "", "UDP Address 127.0.0.1:7789")
	clientCmd.Flags().IntP("max-buffer-len", "l", 1024*10, "Max Buffer Length")
	clientCmd.Flags().IntP("send-peroid-ms", "p", 1, "Send Peroid Ms.")
	clientCmd.Flags().BoolP("recv-check", "c", false, "Recv Check, false is mean to close connection.")
	clientCmd.Flags().StringP("session-id", "s", "bytezero-session-id-0", "Session Id for Channel Create on Tcp.")
	clientCmd.Flags().StringP("device-id", "d", "", "Device Id for app.")
}
