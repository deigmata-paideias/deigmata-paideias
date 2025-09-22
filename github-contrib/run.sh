#!/bin/bash

# GitHub 贡献收集器 v2.0 运行脚本

echo "╔══════════════════════════════════════════════════╗"
echo "║            🚀 GitHub 贡献收集器 v2.0             ║"
echo "║                                                  ║"
echo "║  自动收集和统计您在 GitHub 仓库中的贡献          ║"
echo "╚══════════════════════════════════════════════════╝"
echo

# 检查配置文件
if [ ! -f "config/config.yaml" ]; then
    echo "❌ 配置文件 config.yaml 不存在！"
    echo "请复制 config.example.yaml 为 config.yaml 并填写您的信息。"
    echo
    echo "💡 快速配置:"
    echo "   cp config.example.yaml config.yaml"
    echo "   然后编辑 config.yaml 文件"
    exit 1
fi

echo "✅ 找到配置文件"

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 未找到 Go 环境！请先安装 Go。"
    echo "💡 安装 Go: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go 环境检查通过"

# 检查依赖
echo "📦 检查依赖..."
if ! go mod tidy; then
    echo "❌ 依赖安装失败"
    exit 1
fi

echo "✅ 依赖检查完成"
echo

# 编译程序
echo "🔨 编译程序..."
if go build -o github-contrib ./cmd/github-contrib; then
    echo "✅ 编译成功"
else
    echo "❌ 编译失败"
    echo "💡 请检查Go代码是否有语法错误"
    exit 1
fi

echo

# 检查报告目录
if [ ! -d "reports" ]; then
    echo "📁 创建报告目录..."
    mkdir -p reports
fi

# 运行程序
echo "🚀 开始收集贡献..."
echo "⏱️  这可能需要几分钟，请耐心等待..."
echo

if ./github-contrib; then
    echo
    echo "🎉 收集完成！"
    echo "📄 报告已保存在 reports/ 目录中"
    echo
    echo "📁 生成的文件:"
    ls -la reports/*.md 2>/dev/null || echo "   (暂无报告文件)"
else
    echo
    echo "❌ 执行失败"
    echo "💡 请检查配置文件和网络连接"
    exit 1
fi
