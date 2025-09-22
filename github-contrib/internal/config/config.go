package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置文件结构
type Config struct {
	GitHub struct {
		Username     string   `yaml:"username"`
		Token        string   `yaml:"token"`
		Repositories []string `yaml:"repositories"`
	} `yaml:"github"`
	Output struct {
		ReportDir     string `yaml:"report_dir"`
		IncludeDraft  bool   `yaml:"include_draft"`
		IncludeClosed bool   `yaml:"include_closed"`
	} `yaml:"output"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	var config Config

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证必需的配置
	if config.GitHub.Username == "" {
		return nil, fmt.Errorf("GitHub用户名不能为空")
	}
	if len(config.GitHub.Repositories) == 0 {
		return nil, fmt.Errorf("至少需要配置一个仓库")
	}
	if config.Output.ReportDir == "" {
		config.Output.ReportDir = "./reports"
	}

	return &config, nil
}

// Validate 验证配置有效性
func (c *Config) Validate() error {
	if c.GitHub.Username == "" {
		return fmt.Errorf("GitHub用户名不能为空")
	}

	if len(c.GitHub.Repositories) == 0 {
		return fmt.Errorf("至少需要配置一个仓库")
	}

	for _, repo := range c.GitHub.Repositories {
		if repo == "" {
			return fmt.Errorf("仓库名称不能为空")
		}
	}

	return nil
}
