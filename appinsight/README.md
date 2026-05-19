# AppInsight CLI

面向开发者和 AI Agent 的 iOS App 分析工具。

AppInsight CLI 可以调用 [ipatool](https://github.com/majd/ipatool) 搜索和下载 App Store 的加密 IPA，然后分析 IPA 中可见的信息，包括 Info.plist、权限声明、Framework、资源文件、本地化文案、URL Scheme、Bundle 信息，并生成结构化 JSON 和 Markdown 技术分析报告。

## 重要边界

- 本工具只做开发者分析，不做破解
- 可以调用 ipatool 下载 App Store 加密 IPA
- 不做 IPA 解密
- 不集成 dumpdecrypted、frida-ios-dump、越狱相关能力
- 不绕过 FairPlay / DRM
- 不破解内购
- 不提取用户 token、账号、隐私数据
- 不分发第三方 IPA
- 报告中必须说明：App Store IPA 通常是加密 IPA，无法可靠做代码级反编译，只能基于可见结构、权限、Framework、资源做技术推断

## 安装

### 前置依赖

- Go 1.21+
- macOS（第一版重点支持）
- [ipatool](https://github.com/majd/ipatool)（用于搜索和下载 IPA）
- plutil（macOS 自带）
- strings（Xcode Command Line Tools）

安装 ipatool：

```bash
brew install majd/repo/ipatool
```

安装 Xcode Command Line Tools（如未安装）：

```bash
xcode-select --install
```

### 编译

```bash
git clone <repo-url>
cd appinsight
go build -o bin/appinsight ./cmd/appinsight/
```

## 使用方法

### 1. 环境检查

```bash
appinsight doctor --json
```

输出示例：

```json
{
  "ok": true,
  "command": "doctor",
  "data": {
    "ok": true,
    "command": "doctor",
    "environment": {
      "os": "darwin",
      "arch": "arm64"
    },
    "tools": {
      "ipatool": { "available": true, "path": "/opt/homebrew/bin/ipatool" },
      "plutil": { "available": true, "path": "/usr/bin/plutil" },
      "strings": { "available": true, "path": "/usr/bin/strings" }
    },
    "version": "0.1.0"
  }
}
```

### 2. 搜索 App

```bash
appinsight search "WeChat" --limit 5 --json
```

### 3. 下载 IPA

```bash
appinsight fetch-ios --bundle-id com.tencent.xin --output ./downloads --purchase --json
```

如果 ipatool 未登录，会提示：

```
ipatool is not logged in. Please run: ipatool auth login
```

### 4. 分析 IPA

```bash
appinsight analyze-ipa ./App.ipa --json --output analysis.json
```

分析内容包括：

- Bundle 信息（名称、版本、Build、最低系统版本、设备支持）
- 权限声明（相机、麦克风、定位、通讯录等 13 项权限）
- URL Scheme 和 Query Scheme
- Background Modes
- Framework（系统 Framework + 第三方 SDK 推断）
- 资源文件统计（Asset Catalog、Storyboard、Nib、strings、JSON、Core ML 模型、字体、图片、音频、App Extension）
- 技术栈推断（Flutter、React Native、Unity、Capacitor/Cordova 等）
- LLM 上下文（摘要 + 推荐问题）

### 5. 生成报告

```bash
appinsight report analysis.json --format markdown --output report.md
```

报告包含 9 个章节：

1. 基础信息
2. 分析限制
3. 权限分析
4. Framework 分析
5. 资源结构分析
6. 技术栈推断
7. 可能的实现方式
8. 对开发者的借鉴价值
9. 后续分析建议

## 技术栈推断规则

| 检测特征 | 推断结果 |
|----------|----------|
| Flutter.framework、App.framework、flutter_assets | Flutter |
| React、Hermes、main.jsbundle | React Native |
| UnityFramework.framework | Unity |
| Capacitor、Cordova、www 目录 | Hybrid |
| Vision.framework | 可能使用 Apple Vision |
| CoreImage.framework | 可能做图像处理 |
| CoreML.framework 或 .mlmodelc | 可能使用本地 AI 模型 |
| Metal.framework | 可能使用 GPU 加速 |
| AVFoundation.framework | 可能处理相机、视频、音频 |
| StoreKit.framework | 可能有内购 |
| RevenueCat | 可能使用订阅管理 |
| Firebase | 可能使用统计、远程配置、Crash 或推送 |
| Sentry | 可能使用崩溃监控 |

## 项目结构

```
appinsight/
├── cmd/appinsight/main.go     # 入口
├── internal/
│   ├── cli/                   # cobra 命令定义
│   ├── ipatool/               # ipatool 封装
│   ├── ipa/                   # IPA 分析核心
│   ├── report/                # 报告生成
│   ├── output/                # 统一 JSON 输出
│   └── system/                # 环境检查
├── go.mod
└── README.md
```

## 许可证

MIT
