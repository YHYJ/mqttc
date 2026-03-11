<h1 align="center">Mqttc</h1>
<h3 align="center">MQTT 客户端工具</h3>

<!-- File: README.md -->
<!-- Author: YJ -->
<!-- Email: yj1516268@outlook.com -->
<!-- Created Time: 2023-06-07 11:09:05 -->

---

<p align="center">
  <a href="https://github.com/YHYJ/mqttc/actions/workflows/release.yml"><img src="https://github.com/YHYJ/mqttc/actions/workflows/release.yml/badge.svg" alt="Go build and release by GoReleaser"></a>
</p>

---

## Table of Contents

<!-- vim-markdown-toc GFM -->

* [适配](#适配)
* [安装](#安装)
  * [一键安装](#一键安装)
  * [编译安装](#编译安装)
    * [当前平台](#当前平台)
    * [交叉编译](#交叉编译)
* [用法](#用法)

<!-- vim-markdown-toc -->

---

<!----------------------------------->
<!--                  _   _        -->
<!--  _ __ ___   __ _| |_| |_ ___  -->
<!-- | '_ ` _ \ / _` | __| __/ __| -->
<!-- | | | | | | (_| | |_| || (__  -->
<!-- |_| |_| |_|\__, |\__|\__\___| -->
<!--               |_|             -->
<!----------------------------------->

---

## 适配

- [ ] Linux
- [ ] macOS
- [ ] Windows

## 安装

### 一键安装

```bash
curl -fsSL https://raw.githubusercontent.com/YHYJ/mqttc/main/install.sh | sudo bash -s
```

也可以从 [GitHub Releases](https://github.com/YHYJ/mqttc/releases) 下载解压后使用

### 编译安装

#### 当前平台

如果要为当前平台编译，可以使用以下命令：

```bash
go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/mqttc/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/mqttc/general.BuildTime=`date +%s` -X github.com/yhyj/mqttc/general.BuildBy=$USER" -o build/mqttc main.go
```

#### 交叉编译

> 使用命令`go tool dist list`查看支持的平台
>
> Linux 和 macOS 使用命令`uname -m`，Windows 使用命令`echo %PROCESSOR_ARCHITECTURE%` 确认系统架构
>
> - 例如 x86_64 则设 GOARCH=amd64
> - 例如 aarch64 则设 GOARCH=arm64
> - ...

设置如下系统变量后使用 [编译安装](#编译安装) 的命令即可进行交叉编译：

- CGO_ENABLED: 不使用 CGO，设为 0
- GOOS: 设为 linux, darwin 或 windows
- GOARCH: 根据当前系统架构设置

## 用法

- `config`子命令

  操作配置文件，有以下参数：

  - '--create'：交互式创建配置文件
  - '--open'：使用系统默认编辑器打开配置文件
  - '--print'：打印配置文件内容

- `version`子命令

  查看程序版本信息

- `help`子命令

  查看程序帮助信息
