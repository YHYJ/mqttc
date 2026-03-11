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

	"github.com/pelletier/go-toml"
)

// 一般需要提供的配置
var (
	Host            = "127.0.0.1"      // 默认 MQTT 服务地址
	Port            = 1883             // 默认 MQTT 服务端口
	Username        = ""               // 默认用户名
	Password        = ""               // 默认密码
	PublishTopics   = []string{"test"} // 默认发布主题
	SubscribeTopics = []string{"test"} // 默认订阅主题
)

// 用于转换 Toml 配置树的结构体
type Config struct {
	Host            string   `toml:"host"`
	Port            int      `toml:"port"`
	Username        string   `toml:"username"`
	password        string   `toml:"password"`
	ClientID        string   `toml:"client_id"`
	PublishTopics   []string `toml:"publish_topics"`
	SubscribeTopics []string `toml:"subscribe_topics"`
	QoS             int      `toml:"qos"`
	Retain          bool     `toml:"retain"`
	CleanSession    bool     `toml:"clean_session"`
	Wait            int      `toml:"wait"`
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
	// 一般使用默认值即可的配置
	var (
		ClientID     = ""
		QoS          = 1
		Retain       = false
		CleanSession = true
		Wait         = 5
	)

	// 定义一个 map[string]any 类型的变量并赋值
	exampleConf := map[string]any{
		"host":             Host,            // MQTT 服务地址
		"port":             Port,            // MQTT 服务端口
		"username":         Username,        // 用户名
		"password":         Password,        // 密码
		"client_id":        ClientID,        // 客户端 ID，留空时自动生成
		"publish_topics":   PublishTopics,   // 发布主题,
		"subscribe_topics": SubscribeTopics, // 订阅主题
		"qos":              QoS,             // 服务质量，0/1/2
		"retain":           Retain,          // 是否保留最后一条消息
		"clean_session":    CleanSession,    // 是否清除会话
		"wait":             Wait,            // 连接等待时间
	}
	// 检测配置文件是否存在
	if !FileExist(filePath) {
		return 0, fmt.Errorf("Open %s: no such file or directory", filePath)
	}
	// 检测配置文件是否是 toml 文件
	if !isTomlFile(filePath) {
		return 0, fmt.Errorf("Open %s: is not a toml file", filePath)
	}
	// 把 exampleConf 转换为 *toml.Tree 类型
	tree, err := toml.TreeFromMap(exampleConf)
	if err != nil {
		return 0, err
	}
	// 打开一个文件并获取 io.Writer 接口
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	return tree.WriteTo(file)
}
