/*
File: publish.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2026-03-08 16:41:43

Description: 执行子命令 'publish'
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a message on a topic",
	Long:  `Publish user-input messages on a specified topic.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("publish called")
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
}
