/*
File: root.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2026-03-08 16:06:45

Description: 执行程序
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yhyj/mqttc/general"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mqttc",
	Short: "MQTT client tools",
	Long:  `Simple MQTT client tools.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", general.ConfigFile, "Specify configuration file")

	rootCmd.Flags().BoolP("help", "h", false, "help for mqttc")
}
