#!/bin/bash
set -e

# 项目根目录
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
OUTPUT_DIR="${PROJECT_DIR}/target"
BINARY_NAME="jen"

# 清理旧的构建产物
rm -rf "${OUTPUT_DIR}/${BINARY_NAME}"
mkdir -p "${OUTPUT_DIR}"

# 构建二进制（模板和配置已通过 go:embed 嵌入，无需额外复制资源文件）
echo "🔨 正在构建 ${BINARY_NAME}..."
cd "${PROJECT_DIR}"
go build -o "${OUTPUT_DIR}/${BINARY_NAME}" ./cmd/jen

echo "✅ 构建成功: ${OUTPUT_DIR}/${BINARY_NAME}"
echo ""
echo "使用方式:"
echo "  直接运行:  ${OUTPUT_DIR}/${BINARY_NAME}"
echo "  全局安装:  go install github.com/LingoJack/model_infrax/cmd/jen@latest"