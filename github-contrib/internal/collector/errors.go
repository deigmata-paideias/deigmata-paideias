package collector

import (
	"fmt"
	"strings"
)

// ErrorType 错误类型
type ErrorType int

const (
	ErrorTypeNetwork ErrorType = iota
	ErrorTypeAuth
	ErrorTypeRateLimit
	ErrorTypeNotFound
	ErrorTypePermission
	ErrorTypeConfig
	ErrorTypeUnknown
)

// CollectorError 自定义错误类型
type CollectorError struct {
	Type       ErrorType
	Message    string
	Suggestion string
	Cause      error
}

func (e *CollectorError) Error() string {
	return e.Message
}

// NewCollectorError 创建新的收集器错误
func NewCollectorError(errType ErrorType, message, suggestion string, cause error) *CollectorError {
	return &CollectorError{
		Type:       errType,
		Message:    message,
		Suggestion: suggestion,
		Cause:      cause,
	}
}

// HandleError 处理和格式化错误信息
func HandleError(err error, context string) error {
	if err == nil {
		return nil
	}

	// 检查是否是GitHub API错误
	errMsg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(errMsg, "rate limit"):
		return NewCollectorError(
			ErrorTypeRateLimit,
			fmt.Sprintf("API请求频率限制: %s", context),
			"建议: 1) 配置GitHub Token以提高限制 2) 等待一段时间后重试 3) 减少并发请求",
			err,
		)

	case strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "404"):
		return NewCollectorError(
			ErrorTypeNotFound,
			fmt.Sprintf("仓库不存在或无权限访问: %s", context),
			"建议: 1) 检查仓库名称是否正确 2) 确认仓库是公开的或您有访问权限 3) 检查用户名拼写",
			err,
		)

	case strings.Contains(errMsg, "forbidden") || strings.Contains(errMsg, "403"):
		return NewCollectorError(
			ErrorTypePermission,
			fmt.Sprintf("权限不足: %s", context),
			"建议: 1) 配置有效的GitHub Token 2) 确认Token有足够的权限 3) 检查仓库访问权限",
			err,
		)

	case strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "401"):
		return NewCollectorError(
			ErrorTypeAuth,
			fmt.Sprintf("认证失败: %s", context),
			"建议: 1) 检查GitHub Token是否有效 2) 重新生成Token 3) 确认Token权限范围",
			err,
		)

	case strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "connection"):
		return NewCollectorError(
			ErrorTypeNetwork,
			fmt.Sprintf("网络连接问题: %s", context),
			"建议: 1) 检查网络连接 2) 稍后重试 3) 检查防火墙设置",
			err,
		)

	default:
		return NewCollectorError(
			ErrorTypeUnknown,
			fmt.Sprintf("未知错误: %s - %v", context, err),
			"建议: 1) 检查网络连接 2) 验证配置信息 3) 查看详细错误日志",
			err,
		)
	}
}

// PrintFriendlyError 打印友好的错误信息
func PrintFriendlyError(err error) {
	if collectorErr, ok := err.(*CollectorError); ok {
		fmt.Printf("❌ %s\n", collectorErr.Message)

		if collectorErr.Suggestion != "" {
			fmt.Printf("💡 %s\n", collectorErr.Suggestion)
		}

		if collectorErr.Cause != nil {
			fmt.Printf("🔍 详细错误: %v\n", collectorErr.Cause)
		}
	} else {
		fmt.Printf("❌ 错误: %v\n", err)
	}
}

// IsRetryableError 判断错误是否可以重试
func IsRetryableError(err error) bool {
	if collectorErr, ok := err.(*CollectorError); ok {
		switch collectorErr.Type {
		case ErrorTypeNetwork, ErrorTypeRateLimit:
			return true
		default:
			return false
		}
	}
	return false
}

// GetErrorHelp 获取错误的帮助信息
func GetErrorHelp(errType ErrorType) string {
	switch errType {
	case ErrorTypeAuth:
		return `
GitHub Token 配置帮助:
1. 访问 https://github.com/settings/tokens
2. 点击 "Generate new token (classic)"
3. 选择适当的权限范围:
   - public_repo: 访问公开仓库
   - repo: 访问私有仓库
4. 将生成的token填入config.yaml的token字段
`
	case ErrorTypeRateLimit:
		return `
API限制说明:
- 无认证: 60次/小时
- 有Token: 5000次/小时
- 建议配置GitHub Token以提高限制
`
	case ErrorTypeNotFound:
		return `
仓库访问问题排查:
1. 确认仓库名格式: owner/repo
2. 检查仓库是否存在
3. 确认仓库是公开的或您有访问权限
4. 检查用户名拼写是否正确
`
	default:
		return "请检查网络连接和配置信息"
	}
}
