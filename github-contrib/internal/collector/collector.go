package collector

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"

	"github-contrib/internal/config"
	"github-contrib/internal/template"
)

// GitHubCollector GitHub数据收集器
type GitHubCollector struct {
	client *github.Client
	config *config.Config
}

// ProgressInfo 进度信息
type ProgressInfo struct {
	CurrentRepo    string
	TotalRepos     int
	CurrentRepoIdx int
	CurrentPage    int
	TotalPages     int
	ProcessedPRs   int
	UserPRs        int
}

// NewGitHubCollector 创建新的GitHub收集器
func NewGitHubCollector(cfg *config.Config) (*GitHubCollector, error) {
	var client *github.Client

	fmt.Println("🔧 初始化 GitHub 客户端...")

	if cfg.GitHub.Token != "" {
		fmt.Println("   ✅ 使用 GitHub Token 认证")
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.GitHub.Token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
	} else {
		fmt.Println("   ⚠️  使用无认证模式 (API限制: 60次/小时)")
		client = github.NewClient(nil)
	}

	return &GitHubCollector{
		client: client,
		config: cfg,
	}, nil
}

// printProgress 打印进度信息
func (gc *GitHubCollector) printProgress(info ProgressInfo) {
	fmt.Printf("\r🔍 [%d/%d] %s | 页面: %d | 已处理: %d PR | 找到: %d 个贡献",
		info.CurrentRepoIdx, info.TotalRepos, info.CurrentRepo,
		info.CurrentPage, info.ProcessedPRs, info.UserPRs)
}

// CollectContributions 收集GitHub贡献
func (gc *GitHubCollector) CollectContributions(ctx context.Context, repo string, repoIdx, totalRepos int) (*template.ReportData, error) {
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("仓库格式错误，应为 owner/repo 格式: %s", repo)
	}
	owner, repoName := parts[0], parts[1]

	fmt.Printf("\n📊 [%d/%d] 开始收集 %s 在 %s 的贡献\n", repoIdx, totalRepos, gc.config.GitHub.Username, repo)
	fmt.Println("   📡 连接到 GitHub API...")

	// 首先获取仓库信息和大概的PR数量
	repoInfo, _, err := gc.client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return nil, HandleError(err, fmt.Sprintf("获取仓库 %s 信息", repo))
	}

	fmt.Printf("   📈 仓库信息: %s ⭐%d 🍴%d\n",
		repoInfo.GetDescription(), repoInfo.GetStargazersCount(), repoInfo.GetForksCount())

	var allPRs []*github.PullRequest
	page := 1
	processedPRs := 0
	userPRs := 0

	fmt.Println("   🔍 开始扫描 Pull Requests...")

	// 分页获取所有PR
	for {
		opts := &github.PullRequestListOptions{
			State:     "all",
			Sort:      "updated",
			Direction: "desc",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		}

		progress := ProgressInfo{
			CurrentRepo:    repo,
			TotalRepos:     totalRepos,
			CurrentRepoIdx: repoIdx,
			CurrentPage:    page,
			ProcessedPRs:   processedPRs,
			UserPRs:        userPRs,
		}
		gc.printProgress(progress)

		prs, resp, err := gc.client.PullRequests.List(ctx, owner, repoName, opts)
		if err != nil {
			return nil, HandleError(err, fmt.Sprintf("获取 %s 的PR列表 (页面 %d)", repo, page))
		}

		if len(prs) == 0 {
			break
		}

		// 过滤当前用户的PR
		for _, pr := range prs {
			processedPRs++
			if pr.User != nil && pr.User.GetLogin() == gc.config.GitHub.Username {
				allPRs = append(allPRs, pr)
				userPRs++
			}

			// 更新进度显示
			if processedPRs%10 == 0 {
				progress.ProcessedPRs = processedPRs
				progress.UserPRs = userPRs
				gc.printProgress(progress)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage

		// 避免API限制，稍微延迟
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\n   ✅ 扫描完成! 共处理 %d 个 PR，找到 %d 个您的贡献\n", processedPRs, len(allPRs))

	if len(allPRs) == 0 {
		fmt.Printf("   ℹ️  在仓库 %s 中未找到您的贡献\n", repo)
	}

	fmt.Println("   📝 分析贡献类型...")

	// 分类PR
	var mergedPRs, openPRs, closedPRs []template.PullRequest

	for i, pr := range allPRs {
		if i%5 == 0 {
			fmt.Printf("\r   🔄 分析进度: %d/%d", i+1, len(allPRs))
		}

		prData := template.PullRequest{
			Number:    pr.GetNumber(),
			Title:     pr.GetTitle(),
			URL:       pr.GetHTMLURL(),
			State:     pr.GetState(),
			CreatedAt: pr.GetCreatedAt().Time,
		}

		// 添加标签信息
		if pr.Labels != nil {
			for _, label := range pr.Labels {
				prData.Labels = append(prData.Labels, label.GetName())
			}
		}

		// 根据状态分类
		switch pr.GetState() {
		case "open":
			if gc.config.Output.IncludeDraft || !pr.GetDraft() {
				openPRs = append(openPRs, prData)
			}
		case "closed":
			if pr.GetMergedAt().IsZero() {
				// 已关闭但未合并
				if gc.config.Output.IncludeClosed {
					if pr.ClosedAt != nil {
						prData.ClosedAt = &pr.ClosedAt.Time
					}
					closedPRs = append(closedPRs, prData)
				}
			} else {
				// 已合并
				prData.MergedAt = &pr.MergedAt.Time
				if pr.ClosedAt != nil {
					prData.ClosedAt = &pr.ClosedAt.Time
				}
				mergedPRs = append(mergedPRs, prData)
			}
		}
	}

	reportData := &template.ReportData{
		Username:      gc.config.GitHub.Username,
		Repository:    repo,
		GeneratedAt:   time.Now().Format("2006-01-02 15:04:05"),
		MergedPRs:     mergedPRs,
		OpenPRs:       openPRs,
		ClosedPRs:     closedPRs,
		TotalContribs: len(mergedPRs) + len(openPRs) + len(closedPRs),
	}

	fmt.Printf("\n   📊 贡献统计: ✅已合并 %d | 🔄待处理 %d | ❌已关闭 %d\n",
		len(mergedPRs), len(openPRs), len(closedPRs))

	return reportData, nil
}

// SaveReport 保存报告到文件
func (gc *GitHubCollector) SaveReport(reportData *template.ReportData) error {
	fmt.Printf("💾 生成报告文件...\n")

	// 创建报告目录
	err := os.MkdirAll(gc.config.Output.ReportDir, 0755)
	if err != nil {
		return fmt.Errorf("创建报告目录失败: %w", err)
	}

	fmt.Printf("   📝 渲染 Markdown 模板...\n")

	// 生成报告内容
	content, err := template.GenerateReport(*reportData)
	if err != nil {
		return fmt.Errorf("生成报告失败: %w", err)
	}

	// 生成文件名
	repoName := strings.ReplaceAll(reportData.Repository, "/", "-")
	filename := fmt.Sprintf("%s-%s.md", reportData.Username, repoName)
	filepath := filepath.Join(gc.config.Output.ReportDir, filename)

	fmt.Printf("   📄 保存到文件: %s\n", filename)

	// 写入文件
	err = os.WriteFile(filepath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("保存报告文件失败: %w", err)
	}

	fmt.Printf("   ✅ 报告已保存: %s (%d 字节)\n", filepath, len(content))
	return nil
}

// GetRateLimit 获取API限制信息
func (gc *GitHubCollector) GetRateLimit(ctx context.Context) {
	rateLimit, _, err := gc.client.RateLimits(ctx)
	if err != nil {
		log.Printf("   ⚠️  无法获取API限制信息: %v", err)
		return
	}

	core := rateLimit.GetCore()
	fmt.Printf("   📊 API限制: %d/%d (重置时间: %v)\n",
		core.Remaining, core.Limit, core.Reset.Time.Format("15:04:05"))
}
