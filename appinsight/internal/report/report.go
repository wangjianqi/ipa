package report

import (
	"encoding/json"
	"fmt"
	"html"
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

func GenerateHTML(analysis *ipa.AnalysisResult) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>AppInsight 分析报告 - `)
	sb.WriteString(html.EscapeString(analysis.Bundle.Name))
	sb.WriteString(`</title>
<style>
  :root {
    --bg: #ffffff;
    --surface: #f8f9fa;
    --border: #e1e4e8;
    --text: #24292e;
    --text-secondary: #586069;
    --accent: #0366d6;
    --risk-high: #d73a49;
    --risk-medium: #e36209;
    --risk-low: #28a745;
    --table-header-bg: #f6f8fa;
  }
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    color: var(--text);
    background: var(--bg);
    line-height: 1.6;
    padding: 2rem;
    max-width: 960px;
    margin: 0 auto;
  }
  h1 { font-size: 1.75rem; margin-bottom: 1rem; padding-bottom: 0.5rem; border-bottom: 2px solid var(--accent); }
  h2 { font-size: 1.35rem; margin-top: 2rem; margin-bottom: 0.75rem; padding-bottom: 0.3rem; border-bottom: 1px solid var(--border); }
  h3 { font-size: 1.1rem; margin-top: 1.25rem; margin-bottom: 0.5rem; }
  table { width: 100%; border-collapse: collapse; margin: 0.75rem 0; font-size: 0.9rem; }
  th, td { padding: 0.5rem 0.75rem; border: 1px solid var(--border); text-align: left; }
  th { background: var(--table-header-bg); font-weight: 600; }
  tr:nth-child(even) { background: var(--surface); }
  ul { padding-left: 1.5rem; margin: 0.5rem 0; }
  li { margin: 0.25rem 0; }
  code { background: var(--surface); padding: 0.15rem 0.4rem; border-radius: 4px; font-size: 0.85em; font-family: "SF Mono", Consolas, "Liberation Mono", Menlo, monospace; }
  .risk-high { color: var(--risk-high); font-weight: 600; }
  .risk-medium { color: var(--risk-medium); font-weight: 600; }
  .risk-low { color: var(--risk-low); font-weight: 600; }
  .notice { background: #fff8e1; border-left: 4px solid #f9a825; padding: 0.75rem 1rem; margin: 0.75rem 0; border-radius: 4px; font-size: 0.9rem; }
  .info { background: #e8f4fd; border-left: 4px solid var(--accent); padding: 0.75rem 1rem; margin: 0.75rem 0; border-radius: 4px; font-size: 0.9rem; }
  .footer { margin-top: 2rem; padding-top: 1rem; border-top: 1px solid var(--border); color: var(--text-secondary); font-size: 0.85rem; text-align: center; }
  .tag { display: inline-block; padding: 0.15rem 0.5rem; border-radius: 12px; font-size: 0.8rem; font-weight: 500; margin: 0.1rem 0.2rem; }
  .tag-high { background: #fdecea; color: var(--risk-high); }
  .tag-medium { background: #fff3e0; color: var(--risk-medium); }
  .tag-low { background: #e8f5e9; color: var(--risk-low); }
</style>
</head>
<body>
`)

	sb.WriteString("<h1>AppInsight 分析报告</h1>\n")

	sb.WriteString("<h2>1. 基础信息</h2>\n")
	sb.WriteString("<table>\n")
	sb.WriteString("<tr><th>字段</th><th>值</th></tr>\n")
	sb.WriteString(fmt.Sprintf("<tr><td>App 名称</td><td>%s</td></tr>\n", html.EscapeString(analysis.Bundle.Name)))
	sb.WriteString(fmt.Sprintf("<tr><td>Bundle ID</td><td><code>%s</code></td></tr>\n", html.EscapeString(analysis.Bundle.BundleID)))
	sb.WriteString(fmt.Sprintf("<tr><td>版本</td><td>%s</td></tr>\n", html.EscapeString(analysis.Bundle.Version)))
	sb.WriteString(fmt.Sprintf("<tr><td>Build</td><td>%s</td></tr>\n", html.EscapeString(analysis.Bundle.Build)))
	sb.WriteString(fmt.Sprintf("<tr><td>最低系统版本</td><td>%s</td></tr>\n", html.EscapeString(analysis.Bundle.MinimumOSVersion)))
	sb.WriteString(fmt.Sprintf("<tr><td>支持设备</td><td>%s</td></tr>\n", html.EscapeString(strings.Join(analysis.Bundle.DeviceFamilies, ", "))))
	sb.WriteString("</table>\n")

	sb.WriteString("<h2>2. 分析限制</h2>\n")
	sb.WriteString(fmt.Sprintf("<div class=\"notice\"><strong>加密状态</strong>：%s<br>该 IPA 很可能是 App Store 加密 IPA，无法可靠做代码级反编译和核心算法提取。<br>本报告只基于可见元数据、资源、权限、Framework 和文件结构推断。<br>所有推断均需进一步验证，不应作为确定结论。</div>\n", html.EscapeString(analysis.Encryption.Reason)))

	sb.WriteString("<h2>3. 权限分析</h2>\n")
	if len(analysis.Permissions.Details) == 0 {
		sb.WriteString("<p>未发现声明的敏感权限。</p>\n")
	} else {
		sb.WriteString("<table>\n")
		sb.WriteString("<tr><th>权限 Key</th><th>用途描述</th><th>风险等级</th></tr>\n")
		for _, p := range analysis.Permissions.Details {
			tagClass := "tag-low"
			switch p.Risk {
			case "high":
				tagClass = "tag-high"
			case "medium":
				tagClass = "tag-medium"
			}
			sb.WriteString(fmt.Sprintf("<tr><td><code>%s</code></td><td>%s</td><td><span class=\"tag %s\">%s</span></td></tr>\n",
				html.EscapeString(p.Key), html.EscapeString(p.Description), tagClass, html.EscapeString(p.Risk)))
		}
		sb.WriteString("</table>\n")
	}

	sb.WriteString("<h2>4. Framework 分析</h2>\n")
	sb.WriteString("<h3>系统 Framework</h3>\n")
	if len(analysis.Frameworks.System) == 0 {
		sb.WriteString("<p>未检测到系统 Framework。</p>\n")
	} else {
		sb.WriteString("<ul>\n")
		for _, fw := range analysis.Frameworks.System {
			sb.WriteString(fmt.Sprintf("<li>%s</li>\n", html.EscapeString(fw)))
		}
		sb.WriteString("</ul>\n")
	}

	sb.WriteString("<h3>可能的第三方 SDK</h3>\n")
	if len(analysis.Frameworks.ThirdPartyHints) == 0 {
		sb.WriteString("<p>未检测到第三方 SDK。</p>\n")
	} else {
		sb.WriteString("<ul>\n")
		for _, fw := range analysis.Frameworks.ThirdPartyHints {
			sb.WriteString(fmt.Sprintf("<li>%s</li>\n", html.EscapeString(fw)))
		}
		sb.WriteString("</ul>\n")
	}

	sb.WriteString("<h2>5. 资源结构分析</h2>\n")
	sb.WriteString("<table>\n")
	sb.WriteString("<tr><th>资源类型</th><th>数量</th></tr>\n")
	sb.WriteString(fmt.Sprintf("<tr><td>Asset Catalog (.car)</td><td>%d</td></tr>\n", analysis.Resources.AssetCatalogs))
	sb.WriteString(fmt.Sprintf("<tr><td>Storyboard (.storyboardc)</td><td>%d</td></tr>\n", analysis.Resources.Storyboards))
	sb.WriteString(fmt.Sprintf("<tr><td>Nib (.nib)</td><td>%d</td></tr>\n", analysis.Resources.Nibs))
	sb.WriteString(fmt.Sprintf("<tr><td>Strings 文件</td><td>%d</td></tr>\n", analysis.Resources.StringsFiles))
	sb.WriteString(fmt.Sprintf("<tr><td>JSON 文件</td><td>%d</td></tr>\n", analysis.Resources.JSONFiles))
	sb.WriteString(fmt.Sprintf("<tr><td>Core ML 模型</td><td>%d</td></tr>\n", analysis.Resources.MLModels))
	sb.WriteString(fmt.Sprintf("<tr><td>字体文件</td><td>%d</td></tr>\n", analysis.Resources.Fonts))
	sb.WriteString(fmt.Sprintf("<tr><td>图片文件</td><td>%d</td></tr>\n", analysis.Resources.Images))
	sb.WriteString(fmt.Sprintf("<tr><td>音频文件</td><td>%d</td></tr>\n", analysis.Resources.AudioFiles))
	sb.WriteString(fmt.Sprintf("<tr><td>App Extension</td><td>%d</td></tr>\n", analysis.Resources.AppExtensions))
	sb.WriteString("</table>\n")

	if len(analysis.URLSchemes) > 0 {
		sb.WriteString("<h3>URL Schemes</h3>\n<ul>\n")
		for _, s := range analysis.URLSchemes {
			sb.WriteString(fmt.Sprintf("<li><code>%s</code></li>\n", html.EscapeString(s)))
		}
		sb.WriteString("</ul>\n")
	}

	if len(analysis.QuerySchemes) > 0 {
		sb.WriteString("<h3>Query Schemes</h3>\n<ul>\n")
		for _, s := range analysis.QuerySchemes {
			sb.WriteString(fmt.Sprintf("<li><code>%s</code></li>\n", html.EscapeString(s)))
		}
		sb.WriteString("</ul>\n")
	}

	if len(analysis.BackgroundModes) > 0 {
		sb.WriteString("<h3>Background Modes</h3>\n<ul>\n")
		for _, m := range analysis.BackgroundModes {
			sb.WriteString(fmt.Sprintf("<li><code>%s</code></li>\n", html.EscapeString(m)))
		}
		sb.WriteString("</ul>\n")
	}

	sb.WriteString("<h2>6. 技术栈推断</h2>\n")
	sb.WriteString("<table>\n")
	sb.WriteString(fmt.Sprintf("<tr><td>可能的语言</td><td>%s</td></tr>\n", html.EscapeString(strings.Join(analysis.TechStackInference.PossibleLanguages, ", "))))
	sb.WriteString(fmt.Sprintf("<tr><td>可能的框架</td><td>%s</td></tr>\n", html.EscapeString(strings.Join(analysis.TechStackInference.PossibleFrameworks, ", "))))
	if len(analysis.TechStackInference.PossibleSDKs) > 0 {
		sb.WriteString(fmt.Sprintf("<tr><td>可能的 SDK</td><td>%s</td></tr>\n", html.EscapeString(strings.Join(analysis.TechStackInference.PossibleSDKs, ", "))))
	}
	sb.WriteString("</table>\n")
	sb.WriteString("<div class=\"info\">以上技术栈为基于 Framework 和资源结构的推测，可能不完全准确。需要进一步验证。</div>\n")

	sb.WriteString("<h2>7. 可能的实现方式</h2>\n")
	if len(analysis.TechStackInference.Capabilities) > 0 {
		sb.WriteString("<p>基于权限、Framework 和资源分析，该应用可能具备以下能力：</p>\n<ul>\n")
		for _, cap := range analysis.TechStackInference.Capabilities {
			sb.WriteString(fmt.Sprintf("<li>%s</li>\n", html.EscapeString(cap)))
		}
		sb.WriteString("</ul>\n")
	} else {
		sb.WriteString("<p>基于当前分析，未发现明显的特殊实现方式。</p>\n")
	}

	sb.WriteString("<h2>8. 对开发者的借鉴价值</h2>\n")
	sb.WriteString("<ul>\n")
	sb.WriteString(fmt.Sprintf("<li><strong>功能结构</strong>：%s 声明了 %d 项权限，使用了 %d 个系统 Framework。</li>\n",
		html.EscapeString(analysis.Bundle.Name), len(analysis.Permissions.Details), len(analysis.Frameworks.System)))
	if len(analysis.Frameworks.ThirdPartyHints) > 0 {
		sb.WriteString(fmt.Sprintf("<li><strong>技术路线</strong>：可能使用了 %s 等第三方技术。</li>\n",
			html.EscapeString(strings.Join(analysis.Frameworks.ThirdPartyHints, "、"))))
	}
	if analysis.Resources.MLModels > 0 {
		sb.WriteString("<li><strong>AI 能力</strong>：包含 Core ML 模型，可能集成了本地 AI 推理能力，值得借鉴。</li>\n")
	}
	sb.WriteString("<li><strong>无法确认的部分</strong>：由于 IPA 加密，无法确认具体代码实现、架构模式、网络请求细节等。</li>\n")
	sb.WriteString("</ul>\n")

	sb.WriteString("<h2>9. 后续分析建议</h2>\n")
	sb.WriteString("<ol>\n")
	suggestions := []string{
		"结合 App Store 截图，对照权限和 Framework 推断功能模块。",
		"阅读用户评论，了解核心功能和用户痛点。",
		"查看更新记录，追踪功能迭代方向。",
		"访问官网和文档，了解产品定位和技术博客。",
		"手动体验 App 并录屏，结合分析结果做功能映射。",
		"使用 strings 工具进一步提取二进制中的可读字符串，寻找 API 端点、类名等线索。",
	}
	for _, s := range suggestions {
		sb.WriteString(fmt.Sprintf("<li>%s</li>\n", html.EscapeString(s)))
	}
	sb.WriteString("</ol>\n")

	sb.WriteString("<div class=\"footer\">本报告由 AppInsight CLI 自动生成，仅供开发者技术分析参考。</div>\n")
	sb.WriteString("</body>\n</html>\n")

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
