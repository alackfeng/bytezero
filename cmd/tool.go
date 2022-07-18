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

//	"github.com/alackfeng/bytezero/apps/tool"
	"github.com/alackfeng/bytezero/cores/utils"
	"github.com/spf13/cobra"
)

// toolCmd represents the tool command
var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tool called")

//        tool.NewAppsTool().Main()
body := []byte(`{"data":{"phone":"18011721806","authCode":"123321","biz":"SealMiniCaptcha:login"}}`)
        resp, err := utils.NewHttpClient().Post("http://120.25.216.247:3128/sms/sendSealAuthCode", body, "", "")
        if err != nil {
            fmt.Println("--------err ", err.Error())
        }
        fmt.Println("---------resp: ", string(resp))
	},
}

func init() {
	rootCmd.AddCommand(toolCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// toolCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// toolCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
