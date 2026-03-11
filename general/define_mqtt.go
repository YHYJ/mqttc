/*
File: define_mqtt.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2026-03-11 16:02:01

Description: 和 MQTT 建立连接
*/

package general

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gookit/color"
)

func MqttClient(config *Config) {
	var mode string
	switch Mode {
	case "subscribe":
		mode = "订阅模式"
	case "publish":
		mode = "发布模式"
	default:
		mode = "未知模式"
	}
	color.Printf("运行模式: %s\n", mode)
	color.Printf("连接参数:\n")
	color.Printf("  Host          : %s\n", config.Host)
	color.Printf("  Port          : %d\n", config.Port)
	color.Printf("  Username      : %s\n", config.Username)
	color.Printf("  Client ID     : %s\n", config.ClientID)
	if Mode == "publish" {
		color.Printf("  Topic (pub)   : %s\n", config.PublishTopic)
	} else {
		color.Printf("  Topic (sub)   : %s\n", config.SubscribeTopic)
	}
	color.Printf("  QoS           : %d\n", config.QoS)
	color.Printf("  Timeout       : %d\n", config.Timeout)
	color.Printf("  Retain        : %v\n", config.Retain)
	color.Printf("  Clean Session : %v\n\n", config.CleanSession)

	// 构建 MQTT 客户端选项
	broker := color.Sprintf("tcp://%s:%d", config.Host, config.Port)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetCleanSession(config.CleanSession)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(false)
	opts.SetConnectRetryInterval(time.Duration(config.Timeout) * time.Second)
	opts.SetConnectTimeout(time.Duration(config.Timeout) * time.Second)
	opts.SetMaxReconnectInterval(time.Duration(config.Timeout) * time.Second)

	// 创建一个通道用于同步 OnConnect 完成
	connected := make(chan struct{})

	// 连接成功回调
	opts.OnConnect = func(c mqtt.Client) {
		// 确保函数退出时关闭通道，通知主 goroutine
		defer close(connected)

		color.Println("--- <Info> 已连接到 MQTT 服务器")
		if Mode == "subscribe" {
			if token := c.Subscribe(config.SubscribeTopic, byte(config.QoS), messageHandler); token.Wait() && token.Error() != nil {
				log.Printf("--- <Error> 订阅失败: %v", token.Error())
			} else {
				color.Println("--- <Info> 已订阅主题")
			}
		}
		// 如果不需要订阅，通道依然会关闭
	}

	// 连接丢失回调
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		color.Printf("--- <Warning> 连接丢失: %v\n", err)
		color.Printf("--- <Info> 将在 %d 秒后尝试重连...\n", config.Timeout)
	}

	// 创建客户端
	client := mqtt.NewClient(opts)

	// 尝试连接
	color.Println("--- <Info> 正在连接 MQTT 服务器")
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		// 连接失败，打印错误并退出
		color.Printf("--- <Error> 连接失败: %v\n", token.Error())
		os.Exit(1)
	}

	// 等待 OnConnect 执行完毕（包括订阅完成）
	<-connected

	// 根据模式执行操作
	if Mode == "publish" {
		// 发布模式：读取用户输入并发布
		go publishLoop(client, config)
	} else {
		// 订阅模式：消息已在 OnConnect 中订阅，这里只需等待中断
		color.Println("--- <Info> 等待消息中...")
	}

	// 等待退出信号（Ctrl+C）
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Disconnect(250)
	color.Println("已退出")
}

// publishLoop 循环读取标准输入并发布消息
func publishLoop(client mqtt.Client, config *Config) {
	scanner := bufio.NewScanner(os.Stdin)
	color.Printf("--- <Info> 输入要发布的消息 (按 Enter 发送, Ctrl+C 退出):\n")
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		token := client.Publish(config.PublishTopic, byte(config.QoS), config.Retain, text)
		token.Wait()
		if token.Error() != nil {
			log.Printf("--- <Error> 发布失败: %v", token.Error())
		} else {
			color.Printf("--- <Info> 消息已发布: %s\n", text)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("--- <Error> 读取输入错误: %v", err)
	}
}

// messageHandler 订阅消息的回调函数
func messageHandler(client mqtt.Client, msg mqtt.Message) {
	color.Printf("--- <Info> 收到消息 [%s]: %s\n\n", msg.Topic(), string(msg.Payload()))
}
