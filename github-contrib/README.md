# GitHub 贡献收集器 v2.0

一个功能强大的 Go 工具，用于自动收集和统计您在 GitHub 仓库中的贡献，生成详细的 Markdown 报告。

## ✨ 新特性 (v2.0)

- 🎯 **实时进度显示** - 清晰的进度条和统计信息
- 🏗️ **标准项目结构** - 采用Go标准项目布局
- 🛡️ **智能错误处理** - 友好的错误提示和解决建议
- 📊 **详细统计报告** - 包含处理时间、成功率等信息
- 🔄 **重试机制** - 自动识别可重试错误
- 🎨 **美化输出** - 清晰的界面和状态显示

## 功能特性

- 📊 自动收集指定仓库中的 Pull Request 贡献
- 📝 生成详细的 Markdown 格式报告
- 🏷️ 支持按状态分类（已合并、待处理、已关闭）
- 🔍 显示贡献的详细信息（创建时间、标签等）
- ⚙️ 支持配置文件自定义设置
- 🔐 支持 GitHub Token 认证（提高 API 限制）
- 🚀 实时进度显示和状态反馈

## 快速开始

### 1. 配置设置

复制示例配置文件并编辑：

```bash
cp config.example.yaml config.yaml
```

编辑 `config.yaml` 文件，填写您的信息：

```yaml
github:
  username: "your-github-username"  # 您的 GitHub 用户名
  token: ""                         # GitHub Token (可选但推荐)
  repositories:
    - "apache/hertzbeat"            # 目标仓库
```

### 2. 运行程序

使用提供的脚本（推荐）：

```bash
./run.sh
```

或者手动运行：

```bash
# 安装依赖
go mod tidy

# 编译
go build -o github-contrib ./cmd/github-contrib

# 运行
./github-contrib
```

### 3. 查看报告

报告将保存在 `reports/` 目录中，文件名格式为：`用户名-仓库名.md`

例如：`username-apache-hertzbeat.md`

## 项目结构

```
github-contrib/
├── cmd/
│   └── github-contrib/           # 主程序入口
│       └── main.go
├── internal/                     # 内部包
│   ├── config/                  # 配置管理
│   │   └── config.go
│   ├── collector/               # 数据收集器
│   │   ├── collector.go
│   │   └── errors.go           # 错误处理
│   └── template/               # 报告模板
│       └── template.go
├── reports/                     # 报告输出目录
├── config.yaml                 # 配置文件
├── config.example.yaml         # 配置文件示例
├── run.sh                      # 运行脚本
├── go.mod                      # Go 模块文件
└── README.md                   # 说明文档
```

## 配置说明

### config.yaml

```yaml
github:
  username: "your-github-username"     # 必填：GitHub 用户名
  token: ""                            # 可选：GitHub Personal Access Token
  repositories:                       # 必填：要分析的仓库列表
    - "apache/hertzbeat"
    - "owner/repo"

output:
  report_dir: "./reports"              # 报告保存目录
  include_draft: false                 # 是否包含草稿 PR
  include_closed: true                 # 是否包含已关闭但未合并的 PR
```

### GitHub Token 配置

为了获得更好的体验和避免 API 限制，强烈建议配置 GitHub Token：

1. 访问 [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. 点击 "Generate new token (classic)"
3. 选择适当的权限：
   - `public_repo`: 访问公开仓库
   - `repo`: 访问私有仓库（如需要）
4. 将生成的 token 填入 `config.yaml`

**API 限制对比：**
- 无 Token: 60 次/小时
- 有 Token: 5000 次/小时

## 运行示例

```bash
$ ./run.sh

╔══════════════════════════════════════════════════╗
║            � GitHub 贡献收集器 v2.0             ║
║                                                  ║
║  自动收集和统计您在 GitHub 仓库中的贡献          ║
╚══════════════════════════════════════════════════╝

📋 加载配置文件: config.yaml
✅ 配置加载成功
   👤 用户: your-username
   📁 输出目录: ./reports
   � 目标仓库: 1 个
      1. apache/hertzbeat

🔧 初始化收集器...
   ✅ 使用 GitHub Token 认证

🔍 检查 GitHub API 状态...
   📊 API限制: 4999/5000 (重置时间: 14:30:15)

🚀 开始收集贡献 (共 1 个仓库)
------------------------------------------------------------

📊 [1/1] 开始收集 your-username 在 apache/hertzbeat 的贡献
   📡 连接到 GitHub API...
   📈 仓库信息: HertzBeat(赫兹跳动)是一个拥有强大自定义监控能力 ⭐3421 🍴615
   🔍 开始扫描 Pull Requests...
🔍 [1/1] apache/hertzbeat | 页面: 3 | 已处理: 285 PR | 找到: 5 个贡献
   ✅ 扫描完成! 共处理 285 个 PR，找到 5 个您的贡献
   📝 分析贡献类型...
   🔄 分析进度: 5/5
   📊 贡献统计: ✅已合并 3 | 🔄待处理 1 | ❌已关闭 1

💾 生成报告文件...
   📝 渲染 Markdown 模板...
   📄 保存到文件: your-username-apache-hertzbeat.md
   ✅ 报告已保存: ./reports/your-username-apache-hertzbeat.md (2847 字节)
✅ apache/hertzbeat 处理完成
------------------------------------------------------------

============================================================
📈 收集完成总结
============================================================
📊 仓库统计: 1 个仓库, 1 个成功, 0 个失败
🎯 贡献总数: 5 个
⏱️  处理时间: 12.34 秒

✅ 成功处理的仓库:
   • apache/hertzbeat: 5 个贡献 → ./reports/your-username-apache-hertzbeat.md
============================================================
🎉 所有报告生成完成！
```

## 报告示例

生成的报告包含以下内容：

- 📈 贡献统计总览
- ✅ 已合并的贡献列表
- 🔄 待处理的贡献列表
- ❌ 已关闭的贡献列表（可选）

每个贡献条目包含：
- PR 编号和标题
- GitHub 链接
- 创建时间、合并时间等
- 标签信息

## 错误处理

工具内置了智能错误处理，会针对不同类型的错误提供相应的解决建议：

- **网络错误**: 检查网络连接
- **认证错误**: 检查 GitHub Token
- **权限错误**: 确认仓库访问权限
- **API 限制**: 建议配置 Token 或等待
- **仓库不存在**: 检查仓库名称格式

## 故障排除

### 常见问题

1. **API 限制错误**
   ```
   💡 建议: 1) 配置GitHub Token以提高限制 2) 等待一段时间后重试
   ```

2. **仓库不存在**
   ```
   💡 建议: 1) 检查仓库名称是否正确 2) 确认仓库是公开的或您有访问权限
   ```

3. **编译错误**
   ```bash
   # 检查Go版本
   go version
   
   # 清理并重新安装依赖
   go clean -modcache
   go mod download
   ```

### 调试模式

如需查看详细错误信息，可以直接运行编译后的程序：

```bash
./github-contrib
```

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个工具！

## 更新日志

### v2.0 (当前版本)
- ✨ 全新的项目结构和代码组织
- 🎯 实时进度显示和详细统计
- 🛡️ 智能错误处理和友好提示
- 📊 改进的报告格式和信息展示
- 🔄 自动重试机制
- 🎨 美化的命令行界面

### v1.0
- 🚀 基础的贡献收集功能
- 📝 Markdown 报告生成
- ⚙️ 配置文件支持

## 许可证

本项目使用 MIT 许可证。
