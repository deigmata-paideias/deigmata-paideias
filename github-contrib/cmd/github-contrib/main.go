package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github-contrib/internal/collector"
	"github-contrib/internal/config"
)

func printBanner() {
	fmt.Println(`
╔══════════════════════════════════════════════════╗
║            🚀 GitHub 贡献收集器 v2.0             ║
║                                                  ║
║  自动收集和统计您在 GitHub 仓库中的贡献          ║
╚══════════════════════════════════════════════════╝
`)
}

func printSummary(results []CollectionResult) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📈 收集完成总结")
	fmt.Println(strings.Repeat("=", 60))

	totalRepos := len(results)
	successCount := 0
	totalContribs := 0

	for _, result := range results {
		if result.Success {
			successCount++
			totalContribs += result.ContribCount
		}
	}

	fmt.Printf("📊 仓库统计: %d 个仓库, %d 个成功, %d 个失败\n",
		totalRepos, successCount, totalRepos-successCount)
	fmt.Printf("🎯 贡献总数: %d 个\n", totalContribs)
	fmt.Printf("⏱️  处理时间: %.2f 秒\n", time.Since(startTime).Seconds())

	if successCount > 0 {
		fmt.Println("\n✅ 成功处理的仓库:")
		for _, result := range results {
			if result.Success {
				fmt.Printf("   • %s: %d 个贡献 → %s\n",
					result.Repository, result.ContribCount, result.ReportPath)
			}
		}
	}

	if successCount < totalRepos {
		fmt.Println("\n❌ 处理失败的仓库:")
		for _, result := range results {
			if !result.Success {
				fmt.Printf("   • %s: %s\n", result.Repository, result.Error)
			}
		}
	}

	fmt.Println(strings.Repeat("=", 60))
}

type CollectionResult struct {
	Repository   string
	Success      bool
	ContribCount int
	ReportPath   string
	Error        string
}

var startTime time.Time

func main() {
	startTime = time.Now()
	printBanner()

	configPath := "config/config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	fmt.Printf("📋 加载配置文件: %s\n", configPath)

	// 加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	fmt.Printf("✅ 配置加载成功\n")
	fmt.Printf("   👤 用户: %s\n", cfg.GitHub.Username)
	fmt.Printf("   📁 输出目录: %s\n", cfg.Output.ReportDir)
	fmt.Printf("   📦 目标仓库: %d 个\n", len(cfg.GitHub.Repositories))

	for i, repo := range cfg.GitHub.Repositories {
		fmt.Printf("      %d. %s\n", i+1, repo)
	}

	// 创建收集器
	fmt.Println("\n🔧 初始化收集器...")
	coll, err := collector.NewGitHubCollector(cfg)
	if err != nil {
		log.Fatalf("❌ 初始化失败: %v", err)
	}

	ctx := context.Background()

	// 检查API限制
	fmt.Println("🔍 检查 GitHub API 状态...")
	coll.GetRateLimit(ctx)

	var results []CollectionResult

	fmt.Printf("\n🚀 开始收集贡献 (共 %d 个仓库)\n", len(cfg.GitHub.Repositories))
	fmt.Println(strings.Repeat("-", 60))

	// 为每个仓库生成报告
	for i, repo := range cfg.GitHub.Repositories {
		result := CollectionResult{Repository: repo}

		reportData, err := coll.CollectContributions(ctx, repo, i+1, len(cfg.GitHub.Repositories))
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			collector.PrintFriendlyError(err)

			// 如果是可重试的错误，给出重试建议
			if collector.IsRetryableError(err) {
				fmt.Println("   🔄 这是一个可重试的错误，建议稍后重试")
			}

			results = append(results, result)
			continue
		}

		err = coll.SaveReport(reportData)
		if err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("保存报告失败: %v", err)
			collector.PrintFriendlyError(err)
			results = append(results, result)
			continue
		}

		result.Success = true
		result.ContribCount = reportData.TotalContribs
		result.ReportPath = fmt.Sprintf("%s/%s-%s.md", cfg.Output.ReportDir,
			reportData.Username, strings.ReplaceAll(repo, "/", "-"))

		results = append(results, result)

		fmt.Printf("✅ %s 处理完成\n", repo)

		// 添加延迟避免API限制
		if i < len(cfg.GitHub.Repositories)-1 {
			fmt.Println("   ⏳ 等待 1 秒...")
			time.Sleep(1 * time.Second)
		}

		fmt.Println(strings.Repeat("-", 60))
	}

	printSummary(results)
	fmt.Println("🎉 所有报告生成完成！")
}
