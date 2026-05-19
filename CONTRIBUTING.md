# 贡献指南

感谢您对 AppInsight 项目的关注！我们欢迎并感谢所有形式的贡献。

## 如何贡献

### 报告问题

如果您发现了 bug 或有功能建议，请通过 GitHub Issues 提交：

1. 检查是否已有相关 Issue
2. 如果没有，创建新的 Issue 并详细描述：
   - 问题描述
   - 复现步骤
   - 期望行为
   - 实际行为
   - 环境信息（Go 版本、macOS 版本等）

### 提交代码

1. **Fork 仓库**
   ```bash
   # 点击 GitHub 页面的 Fork 按钮
   ```

2. **克隆您的 Fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/appinsight.git
   cd appinsight
   ```

3. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/your-bug-fix
   ```

4. **进行更改**
   - 编写清晰的代码
   - 遵循 Go 代码规范
   - 添加必要的测试

5. **提交更改**
   ```bash
   git add .
   git commit -m "feat: 添加新功能"
   # 或
   git commit -m "fix: 修复某个问题"
   ```

   提交信息格式：
   - `feat:` 新功能
   - `fix:` 修复问题
   - `docs:` 文档更新
   - `refactor:` 代码重构
   - `test:` 测试相关
   - `chore:` 构建/工具相关

6. **推送到 Fork**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **创建 Pull Request**
   - 在 GitHub 上创建 PR
   - 描述更改内容和原因
   - 关联相关 Issue（如果有）

## 开发规范

### 代码风格

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码
- 使用 `go vet` 静态分析

### 项目结构

```
appinsight/
├── cmd/           # 主程序入口
├── internal/      # 内部包
│   ├── cli/       # CLI 命令
│   ├── ipa/       # IPA 分析
│   ├── ipatool/   # ipatool 封装
│   ├── output/    # 输出处理
│   ├── report/    # 报告生成
│   └── system/    # 系统检查
└── ...
```

### 测试

- 为新功能添加单元测试
- 确保所有测试通过
- 测试覆盖率尽可能高

```bash
# 运行测试
go test ./...

# 运行测试并生成覆盖率报告
go test -cover ./...
```

### 文档

- 更新 README.md（如果需要）
- 为导出的函数添加注释
- 更新 CHANGELOG.md（重大更改）

## 行为准则

### 我们的承诺

为了营造一个开放和友好的环境，我们作为贡献者和维护者承诺：

- 尊重所有参与者
- 接受建设性批评
- 关注对社区最有利的事情
- 对其他社区成员表示同理心

### 不可接受的行为

- 使用带有性暗示的语言或图像
- 挑衅、侮辱或贬低性评论
- 公开或私下的骚扰
- 未经明确许可发布他人的私人信息
- 其他不道德或不专业的行为

## 审查流程

1. 维护者会审查您的 PR
2. 可能需要根据反馈进行修改
3. 通过审查后会被合并

## 发布流程

- 版本号遵循 [Semantic Versioning](https://semver.org/)
- 重大更新会记录在 CHANGELOG.md 中

## 获取帮助

如果您需要帮助：

- 查看 [README.md](README.md)
- 搜索 [Issues](https://github.com/wjq/appinsight/issues)
- 创建新的 Issue 提问

## 致谢

再次感谢您的贡献！
