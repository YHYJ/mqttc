/*
File: publish.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2026-03-08 15:56:53

Description: 子命令 'publish' 的实现
*/

package cli

import (
	"github.com/yhyj/mqttc/general"
)

func SetPublishMode() {
	// 设置运行模式
	general.Mode = "publish"
}
