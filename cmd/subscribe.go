/*
File: subscribe.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2026-03-08 16:41:43

Description: 执行子命令 'subscribe'
*/

package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/yhyj/mqttc/cli"
	"github.com/yhyj/mqttc/general"
)

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to topic",
	Long:  `Subscribe to a set of topics and print all messages it receives.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		configFile, _ := cmd.Flags().GetString("config")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		id, _ := cmd.Flags().GetString("id")
		topic, _ := cmd.Flags().GetString("topic")
		qos, _ := cmd.Flags().GetInt("qos")
		timeout, _ := cmd.Flags().GetInt("timeout")
		retain, _ := cmd.Flags().GetBool("retain")
		clean, _ := cmd.Flags().GetBool("clean")

		// 读取配置文件
		configTree, err := general.GetTomlConfig(configFile)
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
			return
		}
		// 获取配置项
		config, err := general.LoadConfigToStruct(configTree)
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
			return
		}

		// 命令行参数优先
		if host != "" {
			config.Host = host
		}
		if port >= 1 && port <= 65535 {
			config.Port = port
		} else {
			color.Printf("%s %s\n", general.WarnText(general.WarnInfoFlag), general.SecondaryText("Port number range [1, 65535], use default value"))
		}
		if username != "" {
			config.Username = username
		}
		if password != "" {
			config.Password = password
		}
		if id != config.ClientID {
			config.ClientID = id
		}
		if len(topic) > 0 {
			config.SubscribeTopic = topic
		}
		if qos >= 0 && qos <= 2 {
			config.QoS = qos
		} else {
			color.Printf("%s %s\n", general.WarnText(general.WarnInfoFlag), general.SecondaryText("QoS range 0/1/2, use default value"))
		}
		if timeout > 0 {
			config.Timeout = timeout
		}
		if retain {
			config.Retain = retain
		}
		if clean {
			config.CleanSession = clean
		}

		// 校验 Client ID
		if config.ClientID == "" {
			config.ClientID = general.StringGen(8, general.Name, "Subscribe", "Gen")
		}

		cli.SetSubscribeMode()
		general.MqttClient(config)
	},
}

func init() {
	subscribeCmd.Flags().String("host", general.Host, "MQTT broker host")
	subscribeCmd.Flags().Int("port", general.Port, "MQTT broker port")
	subscribeCmd.Flags().String("username", general.Username, "MQTT broker username")
	subscribeCmd.Flags().String("password", general.Password, "MQTT broker password")
	subscribeCmd.Flags().String("id", general.ClientID, "MQTT client ID")
	subscribeCmd.Flags().String("topic", general.SubscribeTopic, "Topic to publish to")
	subscribeCmd.Flags().Int("qos", general.QoS, "QoS level")
	subscribeCmd.Flags().Int("timeout", general.Timeout, "Connection timeout (seconds)")
	subscribeCmd.Flags().Bool("retain", general.Retain, "Retain the message")
	subscribeCmd.Flags().Bool("clean", general.CleanSession, "Clean session")

	subscribeCmd.Flags().BoolP("help", "h", false, "help for subscribe command")

	rootCmd.AddCommand(subscribeCmd)
}
