/*
File: publish.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2026-03-08 16:41:43

Description: 执行子命令 'publish'
*/

package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/yhyj/mqttc/cli"
	"github.com/yhyj/mqttc/general"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a message on a topic",
	Long:  `Publish user-input messages on a specified topic.`,
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
			config.PublishTopic = topic
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
			config.ClientID = general.StringGen(8, general.Name, "Publish", "Gen")
		}

		cli.SetPublishMode()
		general.MqttClient(config)
	},
}

func init() {
	publishCmd.Flags().String("host", general.Host, "MQTT broker host")
	publishCmd.Flags().Int("port", general.Port, "MQTT broker port")
	publishCmd.Flags().String("username", general.Username, "MQTT broker username")
	publishCmd.Flags().String("password", general.Password, "MQTT broker password")
	publishCmd.Flags().String("id", general.ClientID, "MQTT client ID")
	publishCmd.Flags().String("topic", general.PublishTopic, "Topic to publish to")
	publishCmd.Flags().Int("qos", general.QoS, "QoS level")
	publishCmd.Flags().Int("timeout", general.Timeout, "Connection timeout (seconds)")
	publishCmd.Flags().Bool("retain", general.Retain, "Retain the message")
	publishCmd.Flags().Bool("clean", general.CleanSession, "Clean session")

	publishCmd.Flags().BoolP("help", "h", false, "help for publish command")

	rootCmd.AddCommand(publishCmd)
}
