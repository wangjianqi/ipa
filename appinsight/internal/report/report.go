package report

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/wjq/appinsight/internal/ipa"
)

func GenerateMarkdown(analysis *ipa.AnalysisResult) string {
	var sb strings.Builder

	sb.WriteString("# AppInsight 分析报告\n\n")

	sb.WriteString("## 1. 基础信息\n\n")
	sb.WriteString("| 字段 | 值 |\n")
	sb.WriteString("|------|----|\n")
	sb.WriteString(fmt.Sprintf("| App 名称 | %s |\n", analysis.Bundle.Name))
	sb.WriteString(fmt.Sprintf("| Bundle ID | %s |\n", analysis.Bundle.BundleID))
	sb.WriteString(fmt.Sprintf("| 版本 | %s |\n", analysis.Bundle.Version))
	sb.WriteString(fmt.Sprintf("| Build | %s |\n", analysis.Bundle.Build))
	sb.WriteString(fmt.Sprintf("| 最低系统版本 | %s |\n", analysis.Bundle.MinimumOSVersion))
	sb.WriteString(fmt.Sprintf("| 支持设备 | %s |\n", strings.Join(analysis.Bundle.DeviceFamilies, ", ")))
	sb.WriteString("\n")

	sb.WriteString("## 2. 分析限制\n\n")
	sb.WriteString(fmt.Sprintf("- **加密状态**: %s\n", analysis.Encryption.Reason))
	sb.WriteString("- 该 IPA 很可能是 App Store 加密 IPA，无法可靠做代码级反编译和核心算法提取。\n")
	sb.WriteString("- 本报告只基于可见元数据、资源、权限、Framework 和文件结构推断。\n")
	sb.WriteString("- 所有推断均需进一步验证，不应作为确定结论。\n\n")

	sb.WriteString("## 3. 权限分析\n\n")
	if len(analysis.Permissions.Details) == 0 {
		sb.WriteString("未发现声明的敏感权限。\n\n")
	} else {
		sb.WriteString("| 权限 Key | 用途描述 | 风险等级 |\n")
		sb.WriteString("|----------|----------|----------|\n")
		for _, p := range analysis.Permissions.Details {
			sb.WriteString(fmt.Sprintf("| `%s` | %s | %s |\n", p.Key, p.Description, p.Risk))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## 4. Framework 分析\n\n")
	sb.WriteString("### 系统 Framework\n\n")
	if len(analysis.Frameworks.System) == 0 {
		sb.WriteString("未检测到系统 Framework。\n\n")
	} else {
		for _, fw := range analysis.Frameworks.System {
			sb.WriteString(fmt.Sprintf("- %s\n", fw))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("### 可能的第三方 SDK\n\n")
	if len(analysis.Frameworks.ThirdPartyHints) == 0 {
		sb.WriteString("未检测到第三方 SDK。\n\n")
	} else {
		for _, fw := range analysis.Frameworks.ThirdPartyHints {
			sb.WriteString(fmt.Sprintf("- %s\n", fw))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## 5. 资源结构分析\n\n")
	sb.WriteString("| 资源类型 | 数量 |\n")
	sb.WriteString("|----------|------|\n")
	sb.WriteString(fmt.Sprintf("| Asset Catalog (.car) | %d |\n", analysis.Resources.AssetCatalogs))
	sb.WriteString(fmt.Sprintf("| Storyboard (.storyboardc) | %d |\n", analysis.Resources.Storyboards))
	sb.WriteString(fmt.Sprintf("| Nib (.nib) | %d |\n", analysis.Resources.Nibs))
	sb.WriteString(fmt.Sprintf("| Strings 文件 | %d |\n", analysis.Resources.StringsFiles))
	sb.WriteString(fmt.Sprintf("| JSON 文件 | %d |\n", analysis.Resources.JSONFiles))
	sb.WriteString(fmt.Sprintf("| Core ML 模型 | %d |\n", analysis.Resources.MLModels))
	sb.WriteString(fmt.Sprintf("| 字体文件 | %d |\n", analysis.Resources.Fonts))
	sb.WriteString(fmt.Sprintf("| 图片文件 | %d |\n", analysis.Resources.Images))
	sb.WriteString(fmt.Sprintf("| 音频文件 | %d |\n", analysis.Resources.AudioFiles))
	sb.WriteString(fmt.Sprintf("| App Extension | %d |\n", analysis.Resources.AppExtensions))
	sb.WriteString("\n")

	if len(analysis.URLSchemes) > 0 {
		sb.WriteString("### URL Schemes\n\n")
		for _, s := range analysis.URLSchemes {
			sb.WriteString(fmt.Sprintf("- `%s`\n", s))
		}
		sb.WriteString("\n")
	}

	if len(analysis.QuerySchemes) > 0 {
		sb.WriteString("### Query Schemes\n\n")
		for _, s := range analysis.QuerySchemes {
			sb.WriteString(fmt.Sprintf("- `%s`\n", s))
		}
		sb.WriteString("\n")
	}

	if len(analysis.BackgroundModes) > 0 {
		sb.WriteString("### Background Modes\n\n")
		for _, m := range analysis.BackgroundModes {
			sb.WriteString(fmt.Sprintf("- `%s`\n", m))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## 6. 技术栈推断\n\n")
	sb.WriteString(fmt.Sprintf("- **可能的语言**: %s\n", strings.Join(analysis.TechStackInference.PossibleLanguages, ", ")))
	sb.WriteString(fmt.Sprintf("- **可能的框架**: %s\n", strings.Join(analysis.TechStackInference.PossibleFrameworks, ", ")))
	if len(analysis.TechStackInference.PossibleSDKs) > 0 {
		sb.WriteString(fmt.Sprintf("- **可能的 SDK**: %s\n", strings.Join(analysis.TechStackInference.PossibleSDKs, ", ")))
	}
	sb.WriteString("\n")
	sb.WriteString("> **注意**: 以上技术栈为基于 Framework 和资源结构的推测，可能不完全准确。需要进一步验证。\n\n")

	sb.WriteString("## 7. 可能的实现方式\n\n")
	if len(analysis.TechStackInference.Capabilities) > 0 {
		sb.WriteString("基于权限、Framework 和资源分析，该应用可能具备以下能力：\n\n")
		for _, cap := range analysis.TechStackInference.Capabilities {
			sb.WriteString(fmt.Sprintf("- %s\n", cap))
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("基于当前分析，未发现明显的特殊实现方式。\n\n")
	}

	sb.WriteString("## 8. 对开发者的借鉴价值\n\n")
	sb.WriteString(fmt.Sprintf("- **功能结构**: %s 声明了 %d 项权限，使用了 %d 个系统 Framework。\n",
		analysis.Bundle.Name, len(analysis.Permissions.Details), len(analysis.Frameworks.System)))
	if len(analysis.Frameworks.ThirdPartyHints) > 0 {
		sb.WriteString(fmt.Sprintf("- **技术路线**: 可能使用了 %s 等第三方技术。\n",
			strings.Join(analysis.Frameworks.ThirdPartyHints, "、")))
	}
	if analysis.Resources.MLModels > 0 {
		sb.WriteString("- **AI 能力**: 包含 Core ML 模型，可能集成了本地 AI 推理能力，值得借鉴。\n")
	}
	sb.WriteString("- **无法确认的部分**: 由于 IPA 加密，无法确认具体代码实现、架构模式、网络请求细节等。\n\n")

	sb.WriteString("## 9. 后续分析建议\n\n")
	sb.WriteString("1. 结合 App Store 截图，对照权限和 Framework 推断功能模块。\n")
	sb.WriteString("2. 阅读用户评论，了解核心功能和用户痛点。\n")
	sb.WriteString("3. 查看更新记录，追踪功能迭代方向。\n")
	sb.WriteString("4. 访问官网和文档，了解产品定位和技术博客。\n")
	sb.WriteString("5. 手动体验 App 并录屏，结合分析结果做功能映射。\n")
	sb.WriteString("6. 使用 `strings` 工具进一步提取二进制中的可读字符串，寻找 API 端点、类名等线索。\n\n")

	sb.WriteString("---\n\n")
	sb.WriteString("*本报告由 AppInsight CLI 自动生成，仅供开发者技术分析参考。*\n")

	return sb.String()
}

func LoadAnalysis(path string) (*ipa.AnalysisResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read analysis file: %w", err)
	}

	var result ipa.AnalysisResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse analysis JSON: %w", err)
	}

	return &result, nil
}

func WriteReport(content, path string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
