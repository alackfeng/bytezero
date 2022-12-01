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

	"github.com/alackfeng/bytezero/apps/flashsign"
	"github.com/spf13/cobra"
)

// flashsignCmd represents the flashsign command
var flashsignCmd = &cobra.Command{
	Use:   "flashsign",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flashsign called")
		lastReportDate, _ := cmd.Flags().GetString("last-report-date")
		tableField, _ := cmd.Flags().GetString("table-field")
		loop, _ := cmd.Flags().GetBool("loop")
		flashsign.NewFlashSignApp().Main(lastReportDate, tableField, loop)
	},
}

func init() {
	rootCmd.AddCommand(flashsignCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flashsignCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flashsignCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	flashsignCmd.Flags().StringP("last-report-date", "d", "", "t_report_dict.lastReportDate format, eg: 2021-12-08 00:00:00")
	flashsignCmd.Flags().StringP("table-field", "t", "", "cmd: averageAmount30day | revenueMonth ")
	flashsignCmd.Flags().BoolP("loop", "l", false, "loop to now")
}
