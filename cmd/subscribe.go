/*
File: subscribe.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2026-03-08 16:41:43

Description: 执行子命令 'subscribe'
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to topic",
	Long:  `Subscribe to a set of topics and print all messages it receives.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("subscribe called")
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
}
