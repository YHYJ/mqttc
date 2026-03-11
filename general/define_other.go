/*
File: define_other.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2026-03-08 13:35:18

Description: 处理一些杂事
*/

package general

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// StringGen 生成一个指定长度的随机字符串
//
// 参数：
//   - length: 长度
//   - name: 名称
//   - prefix: 前缀
//   - suffix: 后缀
//
// 返回：
//   - 字符串
func StringGen(length int, name string, prefix string, suffix string) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// 如果随机失败，回退到基于时间戳生成
		return color.Sprintf("%s-%s-%d-%s", name, prefix, time.Now().UnixNano(), suffix)
	}
	return name + "-" + prefix + "-" + hex.EncodeToString(b) + "-" + suffix
}

// AreYouSure 获取用户二次确认
//
// 参数：
//   - question: 问题
//   - defaultAnswer: 默认回答，true 或 false
//
// 返回：
//   - true/false
//   - 错误信息
func AreYouSure(question string, defaultAnswer bool) (bool, error) {
	var (
		viewAnswers []string                                 // 显示用可选答案
		answersMap  = map[string]bool{"y": true, "n": false} // 可选答案和实际返回值的映射
		reader      = bufio.NewReader(os.Stdin)              // 标准输入
	)

	// 根据 defaultAnswer 设置显示用的可选答案
	if defaultAnswer == true {
		viewAnswers = []string{"Y", "n"}
	} else {
		viewAnswers = []string{"y", "N"}
	}
	viewAnswersConsortium := strings.Join(viewAnswers, "/") // 显示用可选答案的组合体

	for {
		// 输出问题
		color.Printf("%s %s: ", question, SecondaryText("(", viewAnswersConsortium, ")"))

		// 从标准输入中读取用户的回答
		userRawAnswer, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		// 去除用户回答中的换行符
		userAnswer := strings.TrimSpace(strings.TrimSuffix(userRawAnswer, "\n"))

		// 检查用户回答是否符合要求
		for answer, result := range answersMap {
			if strings.EqualFold(userAnswer, answer) {
				return result, nil
			} else if userAnswer == "" { // 如果用户回答为空，返回默认回答
				return defaultAnswer, nil
			}
		}
	}
}

// GiveYourChoice 给出可选项，获取用户选择
//
// 参数：
//   - tips: 提示信息
//   - options: 可选项
//   - defaultOption: 默认选项的下标（从0开始）
//
// 返回：
//   - 用户的选择（去掉了首尾空格和最后的换行符）
//   - 错误信息
func GiveYourChoice(tips string, options []string, defaultOption int) (string, error) {
	var (
		viewOptions = make([]string, len(options)) // 显示用可选项
	)
	copy(viewOptions, options)

	// defaultOptionIndex 所指定的默认选项在 options 中，将显示用的可选项中的默认选项转换为首字母大写
	if defaultOption >= 0 && defaultOption <= len(viewOptions)-1 {
		titleCase := cases.Title(language.English) // 创建一个 Title 风格的转换器，基于英语规则

		// 将第一个单词的首字母大写
		words := strings.Fields(viewOptions[defaultOption])
		if len(words) > 0 {
			viewOptions[defaultOption] = titleCase.String(words[0])
		}
	} else {
		return "", errors.New("The 'defaultOptionIndex' out of range")
	}
	viewOptionsConsortium := strings.Join(viewOptions, "/") // 显示用的可选项的组合体

	for {
		// 输出问题
		color.Printf("%s %s: ", tips, SecondaryText("(", viewOptionsConsortium, ")"))

		// 从标准输入中读取用户的回答
		reader := bufio.NewReader(os.Stdin)
		userRawChoice, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		// 去除用户回答中的换行符
		userChoice := strings.TrimSpace(strings.TrimSuffix(userRawChoice, "\n"))

		// 检查输入是否在有效答案列表中
		for _, answer := range options {
			if strings.EqualFold(userChoice, answer) {
				return answer, nil
			} else if userChoice == "" { // 如果用户未选择，返回默认选项
				return options[defaultOption], nil
			}
		}
	}
}

// GetUserInputString 获取用户输入的字符串
//
// 参数：
//   - tips: 提示信息
//   - defaultValue: 用户未输入时的默认值
//
// 返回：
//   - 用户输入（去掉了最后的换行符）
//   - 错误信息
func GetUserInputString(tips string, defaultValue string) (string, error) {
	color.Printf("%s %s: ", tips, SecondaryText("(", defaultValue, ")"))

	// 从标准输入中读取用户的回答
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultValue, err
	}

	value := func() string {
		if len(input) <= 1 {
			return defaultValue
		}
		return strings.TrimSuffix(input, "\n")
	}()

	return value, nil
}

// GetUserInputInt 获取用户输入的整数
//
// 参数：
//   - tips: 提示信息
//   - defaultValue: 用户未输入时的默认值
//
// 返回：
//   - 用户输入整数
//   - 错误信息
func GetUserInputInt(tips string, defaultValue int) (int, error) {
	color.Printf("%s %s: ", tips, SecondaryText("(", defaultValue, ")"))

	// 从标准输入中读取用户的回答
	reader := bufio.NewReader(os.Stdin)
	originalValue, err := reader.ReadString('\n')
	if err != nil {
		return defaultValue, err
	}

	// 去除换行和首尾空格
	input := strings.TrimSpace(originalValue)

	// 如果用户直接回车，使用默认值
	if input == "" {
		return defaultValue, nil
	}

	// 尝试将输入转换为整数
	value, err := strconv.Atoi(input)
	if err != nil {
		// 转换失败，返回自定义错误
		return 0, fmt.Errorf("输入无效，请输入一个整数: %v", err)
	}

	return value, nil
}
