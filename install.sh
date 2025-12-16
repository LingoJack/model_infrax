#!/bin/bash
echo "installing..."

# 下载压缩包
curl -LO https://github.com/LingoJack/model_infrax/raw/refs/heads/main/release/model_infrax.tar.gz

# 检查下载是否成功
if [ ! -f model_infrax.tar.gz ]; then
    echo "Error: Failed to download model_infrax.tar.gz"
    exit 1
fi

# 解压到 /Applications/ 目录（需要 sudo 权限）
echo "Extracting files to /Applications/..."
# 重定向 stdin 到 /dev/tty 以便 sudo 可以读取密码
sudo tar -xzf model_infrax.tar.gz -C /Applications/ < /dev/tty

# 检查解压是否成功
if [ ! -d /Applications/model_infrax ]; then
    echo "Error: Failed to extract files"
    rm -rf model_infrax.tar.gz
    exit 1
fi

rm -rf model_infrax.tar.gz
echo "install success"

echo "setting permissions..."
# 设置可执行权限（重定向 stdin 到 /dev/tty 以便 sudo 可以读取密码）
sudo chmod +x /Applications/model_infrax/jen /Applications/model_infrax/jcode < /dev/tty
echo "permissions set success"

echo "setting env..."
# 检测当前使用的 shell
SHELL_RC=""
if [ -n "$ZSH_VERSION" ]; then
    SHELL_RC="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_RC="$HOME/.bashrc"
else
    # 默认尝试 zsh
    SHELL_RC="$HOME/.zshrc"
fi

# 检查 PATH 是否已经包含该目录
if ! grep -q "/Applications/model_infrax" "$SHELL_RC" 2>/dev/null; then
    echo "export PATH=\"/Applications/model_infrax:\$PATH\"" >> "$SHELL_RC"
    echo "Added to $SHELL_RC"
else
    echo "PATH already configured in $SHELL_RC"
fi

# 临时设置当前会话的 PATH
export PATH="/Applications/model_infrax:$PATH"
echo "env set success"

# 验证安装
if command -v jen &> /dev/null; then
    echo "jen version: $(jen -v)"
else
    echo "Installation completed. Please restart your terminal or run: source $SHELL_RC"
fi