/*
File: define_toml.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2026-03-08 15:15:12

Description: 操作 TOML 配置文件
*/

package general

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/pelletier/go-toml"
)

// 用于转换 Toml 配置树的结构体
type Config struct {
	Host           string `toml:"host"`            // MQTT 服务地址
	Port           int    `toml:"port"`            // MQTT 服务端口
	Username       string `toml:"username"`        // 用户名
	Password       string `toml:"password"`        // 密码
	ClientID       string `toml:"client_id"`       // 客户端 ID，留空时自动生成随机 ID
	PublishTopic   string `toml:"publish_topic"`   // 发布主题
	SubscribeTopic string `toml:"subscribe_topic"` // 订阅主题
	QoS            int    `toml:"qos"`             // 服务质量，0/1/2
	Timeout        int    `toml:"timeout"`         // 连接超时时间
	Retain         bool   `toml:"retain"`          // 是否保留最后一条消息
	CleanSession   bool   `toml:"clean_session"`   // 是否清空会话
}

// 配置项
var (
	// 允许用户修改的配置项
	Host           = "127.0.0.1"
	Port           = 1883
	Username       = ""
	Password       = ""
	PublishTopic   = "test"
	SubscribeTopic = "test"
	// 使用默认值的配置项（这里保持首字母大写是因为命令行参数需要）
	ClientID     = ""
	QoS          = 1
	Timeout      = 5
	Retain       = false
	CleanSession = true
)

// 配置
var appConfig = Config{
	Host:           Host,
	Port:           Port,
	Username:       Username,
	Password:       Password,
	ClientID:       ClientID,
	PublishTopic:   PublishTopic,
	SubscribeTopic: SubscribeTopic,
	QoS:            QoS,
	Timeout:        Timeout,
	Retain:         Retain,
	CleanSession:   CleanSession,
}

// isTomlFile 检测文件是不是 toml 文件
//
// 参数：
//   - filePath: 待检测文件路径
//
// 返回：
//   - 是 toml 文件返回 true，否则返回 false
func isTomlFile(filePath string) bool {
	if strings.HasSuffix(filePath, ".toml") {
		return true
	}
	return false
}

// GetTomlConfig 读取 toml 配置文件
//
// 参数：
//   - filePath: toml 配置文件路径
//
// 返回：
//   - toml 配置树
//   - 错误信息
func GetTomlConfig(filePath string) (*toml.Tree, error) {
	if !FileExist(filePath) {
		return nil, fmt.Errorf("Open %s: no such file or directory", filePath)
	}
	if !isTomlFile(filePath) {
		return nil, fmt.Errorf("Open %s: is not a toml file", filePath)
	}
	tree, err := toml.LoadFile(filePath)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// LoadConfigToStruct 将 Toml 配置树加载到结构体
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
//
// 返回：
//   - 结构体
//   - 错误信息
func LoadConfigToStruct(configTree *toml.Tree) (*Config, error) {
	var config Config
	if err := configTree.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// WriteTomlConfig 写入 toml 配置文件
//
// 参数：
//   - filePath: toml 配置文件路径
//
// 返回：
//   - 写入的字节数
//   - 错误信息
func WriteTomlConfig(filePath string) (int64, error) {
	// 打开配置文件
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 写入注释
	manual := color.Sprintf("##\n## %s - %s\n## Generaled on %s\n##\n\n", Name, Version, time.Now().Format("2006-01-02 15:04:05"))
	n, err := file.WriteString(manual)
	if err != nil {
		return int64(n), err
	}

	// 创建编码器并设置顺序保留
	encoder := toml.NewEncoder(file)
	encoder.Order(toml.OrderPreserve)

	if err := encoder.Encode(appConfig); err != nil {
		return int64(n), err
	}

	stat, _ := file.Stat()
	return stat.Size(), nil
}
