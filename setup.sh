#!/bin/bash

# 🚀 FPUSDT - TRON API Go版本 环境配置脚本
# 适用于 Linux 和 macOS 系统

echo "🚀 开始配置 FPUSDT 开发环境..."

# 检测shell类型
if [[ "$SHELL" == *"zsh"* ]]; then
    SHELL_RC="$HOME/.zshrc"
    SHELL_NAME="zsh"
elif [[ "$SHELL" == *"bash"* ]]; then
    SHELL_RC="$HOME/.bashrc"
    SHELL_NAME="bash"
else
    SHELL_RC="$HOME/.profile"
    SHELL_NAME="shell"
fi

echo "📋 检测到您使用的是: $SHELL_NAME"
echo "📄 配置文件路径: $SHELL_RC"

# 检查 Go 是否已安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未检测到 Go 环境"
    echo "请先安装 Go 1.19+ 版本: https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ 检测到 Go 版本: $GO_VERSION"

# 设置 Go 代理环境变量
echo ""
echo "🔧 配置 Go 模块代理..."

# 临时设置（当前会话）
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn

echo "✅ 临时环境变量已设置"

# 永久设置（写入配置文件）
echo ""
echo "📝 写入永久配置到 $SHELL_RC ..."

# 备份原配置文件
if [[ -f "$SHELL_RC" ]]; then
    cp "$SHELL_RC" "$SHELL_RC.backup.$(date +%Y%m%d_%H%M%S)"
    echo "📋 已备份原配置文件"
fi

# 检查是否已存在配置
if grep -q "GO111MODULE" "$SHELL_RC" 2>/dev/null; then
    echo "⚠️  检测到已存在 Go 配置，跳过写入"
else
    cat >> "$SHELL_RC" << 'EOF'

# 🚀 FPUSDT Go 环境配置
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn
EOF
    echo "✅ 永久配置已写入"
fi

# 验证配置
echo ""
echo "🔍 验证当前配置..."
echo "GO111MODULE: $(go env GO111MODULE)"
echo "GOPROXY: $(go env GOPROXY)"
echo "GOSUMDB: $(go env GOSUMDB)"

# 测试网络连接
echo ""
echo "🌐 测试代理连接..."
if go list -m golang.org/x/tools >/dev/null 2>&1; then
    echo "✅ 代理连接正常"
else
    echo "⚠️  代理连接测试失败，可能需要检查网络"
fi

# 安装依赖
echo ""
echo "📦 开始安装项目依赖..."
if go mod tidy; then
    echo "✅ 依赖安装成功"
else
    echo "❌ 依赖安装失败，请检查网络连接或手动执行: go mod tidy"
    exit 1
fi

echo ""
echo "🎉 环境配置完成！"
echo ""
echo "📝 下一步操作："
echo "   1. 重新加载配置: source $SHELL_RC"
echo "   2. 启动服务: go run main.go"
echo "   3. 访问主页: http://localhost:9527"
echo ""
echo "🆘 如有问题，请联系: https://t.me/king_orz"
