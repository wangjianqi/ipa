# AppInsight

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS-lightgrey.svg)](https://www.apple.com/macos)

AppInsight 是一个面向开发者的 iOS 应用分析工具，帮助开发者快速了解 App Store 应用的技术架构、权限声明、框架依赖等信息。

## 功能特性

- 🔍 **App Store 搜索** - 通过关键词搜索 App Store 应用
- 📥 **IPA 下载** - 使用 ipatool 从 App Store 下载应用
- 🔬 **IPA 分析** - 深度分析 IPA 文件的元数据、权限、框架和资源
- 📊 **报告生成** - 生成 Markdown 或 HTML 格式的分析报告
- 🔧 **环境检查** - 检查系统依赖工具是否就绪

## 安装

### 前置要求

- macOS 系统
- [Go 1.26+](https://golang.org/dl/)
- [ipatool](https://github.com/majd/ipatool) - 用于从 App Store 下载应用

```bash
# 安装 ipatool
brew install majd/repo/ipatool

# 登录 ipatool（需要 Apple ID）
ipatool auth login
```

### 通过 Go Install 安装（推荐）

如果你已经安装了 Go 1.26+，可以直接使用 `go install`：

```bash
go install github.com/wjq/appinsight/cmd/appinsight@latest
```

安装后，确保 `$GOPATH/bin` 或 `$HOME/go/bin` 在你的 PATH 中。

### 从源码编译安装

```bash
git clone https://github.com/wangjianqi/ipa.git
cd appinsight/appinsight
go build -o appinsight ./cmd/appinsight

# 可选：安装到系统路径
sudo mv appinsight /usr/local/bin/
```

## 使用方法

### 环境检查

```bash
appinsight doctor
```

检查系统环境是否满足运行要求。

### 搜索应用

```bash
# 搜索应用
appinsight search WeChat

# 限制结果数量
appinsight search "支付宝" --limit 5
```

### 下载 IPA

```bash
# 下载应用（需要 bundle ID）
appinsight fetch-ios --bundle-id com.tencent.xin

# 指定输出目录
appinsight fetch-ios --bundle-id com.tencent.xin --output ./downloads

# 购买应用（如果未购买）
appinsight fetch-ios --bundle-id com.example.app --purchase
```

### 分析 IPA

```bash
# 分析 IPA 文件
appinsight analyze-ipa ./WeChat.ipa

# 输出分析结果到文件
appinsight analyze-ipa ./WeChat.ipa --output analysis.json
```

### 生成报告

```bash
# 从分析结果生成 Markdown 报告
appinsight report analysis.json --format markdown

# 输出 Markdown 报告到文件
appinsight report analysis.json --format markdown --output report.md

# 生成 HTML 报告
appinsight report analysis.json --format html

# 输出 HTML 报告到文件
appinsight report analysis.json --format html --output report.html
```

## 分析内容

AppInsight 可以分析以下信息：

### 基础信息
- App 名称、Bundle ID、版本号
- 支持的设备类型（iPhone/iPad）
- 最低系统版本要求

### 权限分析
- 相机、麦克风、相册访问权限
- 定位、通讯录、日历权限
- 用户追踪、语音识别、Face ID 等

### 框架分析
- 系统 Framework 使用情况
- 第三方 SDK 识别（Flutter、React Native、Unity、Firebase 等）

### 资源结构
- Asset Catalogs、Storyboards、Nibs
- Core ML 模型、字体、图片、音频文件
- App Extensions

### 技术栈推断
- 可能的开发语言和框架
- 应用能力推测（AR、地图、内购等）

## 输出格式

所有命令默认输出 JSON 格式，便于程序化处理：

```json
{
  "ok": true,
  "command": "analyze-ipa",
  "data": {
    "bundle": {
      "name": "Example App",
      "bundleId": "com.example.app",
      "version": "1.0.0"
    },
    ...
  }
}
```

## 项目结构

```
appinsight/
├── cmd/appinsight/         # 主程序入口
├── internal/
│   ├── cli/               # CLI 命令定义
│   ├── ipa/               # IPA 分析逻辑
│   ├── ipatool/           # ipatool 封装
│   ├── output/            # 输出格式化
│   ├── report/            # 报告生成
│   └── system/            # 系统检查
├── go.mod
└── go.sum
```

## 技术栈

- [Go](https://golang.org/) - 编程语言
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [ipatool](https://github.com/majd/ipatool) - App Store 交互

## 贡献指南

欢迎提交 Issue 和 Pull Request！请阅读 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详情。

## 行为准则

本项目遵循 [Contributor Covenant](https://www.contributor-covenant.org/) 行为准则，详情请见 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)。

## 许可证

本项目采用 [MIT License](LICENSE) 开源许可证。

## 免责声明

本工具仅供开发者技术分析和学习使用。使用本工具下载和分析应用时，请遵守相关法律法规和 Apple 的使用条款。

- 下载的 IPA 文件仅供个人技术分析，不得用于商业用途
- 分析结果仅基于可见元数据，不代表应用的真实实现
- 请尊重应用开发者的知识产权

## 致谢

- [ipatool](https://github.com/majd/ipatool) - 提供 App Store 交互能力

## 联系方式

如有问题或建议，欢迎通过以下方式联系：

- 提交 [GitHub Issue](https://github.com/wangjianqi/ipa/issues)
- 发送邮件至 [wangjianqi@aliyun.com](mailto:wangjianqi@aliyun.com)

---

<p align="center">Made with ❤️ by WJQ</p>
