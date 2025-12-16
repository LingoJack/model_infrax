#!/bin/bash

set -e  # 遇到错误时退出

# 函数：检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 函数：验证文件是否有效
verify_file() {
    local file=$1
    if [ ! -f "$file" ] || [ ! -s "$file" ]; then
        return 1
    fi
    return 0
}

echo "============================================"
echo "Model Infrax Installation"
echo "============================================"

# 临时文件路径
TEMP_DIR="/tmp/model_infrax_install"
TAR_FILE="$TEMP_DIR/model_infrax.tar.gz"
INSTALL_DIR="/Applications/model_infrax"
DOWNLOAD_URL="https://github.com/LingoJack/model_infrax/raw/refs/heads/main/release/model_infrax.tar.gz"

# 创建临时目录
echo "Creating temporary directory..."
mkdir -p "$TEMP_DIR"

# 下载压缩包（下载到文件而非管道）
echo "Downloading model_infrax..."
if curl -L --progress-bar -o "$TAR_FILE" "$DOWNLOAD_URL"; then
    echo "Download completed"
else
    echo "Error: Failed to download model_infrax.tar.gz"
    rm -rf "$TEMP_DIR"
    exit 1
fi

# 验证下载的文件
echo "Verifying downloaded file..."
if ! verify_file "$TAR_FILE"; then
    echo "Error: Downloaded file is invalid or empty"
    rm -rf "$TEMP_DIR"
    exit 1
fi
echo "File verification passed"

# 解压到 /Applications/ 目录（需要 sudo 权限）
echo "Extracting files to /Applications/..."
echo "Note: You may be prompted for your password"
if sudo tar -xzf "$TAR_FILE" -C /Applications/; then
    echo "Extraction completed"
else
    echo "Error: Failed to extract files"
    rm -rf "$TEMP_DIR"
    exit 1
fi

# 检查解压是否成功
if [ ! -d "$INSTALL_DIR" ]; then
    echo "Error: Installation directory not found after extraction"
    rm -rf "$TEMP_DIR"
    exit 1
fi

# 清理临时文件
echo "Cleaning up temporary files..."
rm -rf "$TEMP_DIR"

# 设置可执行权限
echo "Setting executable permissions..."
if sudo chmod +x "$INSTALL_DIR/jen" "$INSTALL_DIR/jcode"; then
    echo "Permissions set successfully"
else
    echo "Error: Failed to set executable permissions"
    exit 1
fi

# 验证可执行文件
echo "Verifying executables..."
if [ ! -x "$INSTALL_DIR/jen" ]; then
    echo "Error: jen is not executable"
    exit 1
fi
if [ ! -x "$INSTALL_DIR/jcode" ]; then
    echo "Error: jcode is not executable"
    exit 1
fi
echo "Executables verified"

# 检测当前使用的 shell
echo "Configuring environment variables..."
SHELL_RC=""
if [ -n "$ZSH_VERSION" ]; then
    SHELL_RC="$HOME/.zshrc"
    echo "Detected shell: zsh"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_RC="$HOME/.bashrc"
    echo "Detected shell: bash"
else
    # 默认尝试 zsh（macOS 默认）
    SHELL_RC="$HOME/.zshrc"
    echo "Using default shell config: zsh"
fi

# 检查 PATH 是否已经包含该目录（避免重复添加）
PATH_EXPORT="export PATH=\"$INSTALL_DIR:\$PATH\""
if ! grep -q "$INSTALL_DIR" "$SHELL_RC" 2>/dev/null; then
    echo "$PATH_EXPORT" >> "$SHELL_RC"
    echo "Added to $SHELL_RC"
else
    echo "PATH already configured in $SHELL_RC"
fi

# 临时设置当前会话的 PATH
export PATH="$INSTALL_DIR:$PATH"

# 验证安装
echo "Verifying installation..."
if command_exists jen && jen -v >/dev/null 2>&1; then
    echo "✓ jen installed successfully"
    echo "  Version: $(jen -v)"
else
    echo "⚠ jen command not immediately available"
    echo "  Please restart your terminal or run: source $SHELL_RC"
fi

if command_exists jcode; then
    echo "✓ jcode installed successfully"
else
    echo "⚠ jcode command not immediately available"
    echo "  Please restart your terminal or run: source $SHELL_RC"
fi

echo ""
echo "============================================"
echo "Installation completed successfully!"
echo ""
echo "Installation directory: $INSTALL_DIR"
echo "Configuration file: $SHELL_RC"
echo ""
echo "Next steps:"
echo "  1. Restart your terminal or run: source $SHELL_RC"
echo "  2. Run 'jcode' to open your working directory"
echo "  3. Run 'jen' to generate code"
echo ""
echo "Usage:"
echo "  - Add table structure in schema.sql"
echo "  - Add database config in application.yml"
echo "============================================"
