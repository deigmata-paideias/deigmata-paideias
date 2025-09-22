package template

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

// ReportData 包含报告生成所需的所有数据
type ReportData struct {
	Username      string
	Repository    string
	GeneratedAt   string
	MergedPRs     []PullRequest
	OpenPRs       []PullRequest
	ClosedPRs     []PullRequest
	TotalContribs int
}

// PullRequest 表示一个PR的信息
type PullRequest struct {
	Number    int
	Title     string
	URL       string
	State     string
	CreatedAt time.Time
	MergedAt  *time.Time
	ClosedAt  *time.Time
	Labels    []string
}

// 报告模板
const reportTemplate = `# {{.Repository}} 贡献报告

**用户:** {{.Username}}  
**生成时间:** {{.GeneratedAt}}  
**总贡献数:** {{.TotalContribs}}

---

## 📈 贡献统计

- **已合并 PR:** {{len .MergedPRs}} 个
- **待处理 PR:** {{len .OpenPRs}} 个
{{- if .ClosedPRs}}
- **已关闭 PR:** {{len .ClosedPRs}} 个
{{- end}}

---

{{- if .MergedPRs}}

## ✅ 已合并贡献 ({{len .MergedPRs}} 个)

{{- range .MergedPRs}}
- [#{{.Number}}]({{.URL}}) {{.Title}}
  - 创建时间: {{.CreatedAt.Format "2006-01-02 15:04:05"}}
  {{- if .MergedAt}}
  - 合并时间: {{.MergedAt.Format "2006-01-02 15:04:05"}}
  {{- end}}
  {{- if .Labels}}
  - 标签: {{range $i, $label := .Labels}}{{if $i}}, {{end}}` + "`{{$label}}`" + `{{end}}
  {{- end}}
{{- end}}

{{- end}}

{{- if .OpenPRs}}

## 🔄 待处理贡献 ({{len .OpenPRs}} 个)

{{- range .OpenPRs}}
- [#{{.Number}}]({{.URL}}) {{.Title}}
  - 创建时间: {{.CreatedAt.Format "2006-01-02 15:04:05"}}
  {{- if .Labels}}
  - 标签: {{range $i, $label := .Labels}}{{if $i}}, {{end}}` + "`{{$label}}`" + `{{end}}
  {{- end}}
{{- end}}

{{- end}}

{{- if .ClosedPRs}}

## ❌ 已关闭贡献 ({{len .ClosedPRs}} 个)

{{- range .ClosedPRs}}
- [#{{.Number}}]({{.URL}}) {{.Title}}
  - 创建时间: {{.CreatedAt.Format "2006-01-02 15:04:05"}}
  {{- if .ClosedAt}}
  - 关闭时间: {{.ClosedAt.Format "2006-01-02 15:04:05"}}
  {{- end}}
  {{- if .Labels}}
  - 标签: {{range $i, $label := .Labels}}{{if $i}}, {{end}}` + "`{{$label}}`" + `{{end}}
  {{- end}}
{{- end}}

{{- end}}

---

*此报告由 GitHub 贡献收集器自动生成*
`

// GenerateReport 生成markdown报告
func GenerateReport(data ReportData) (string, error) {
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("执行模板失败: %w", err)
	}

	return buf.String(), nil
}
