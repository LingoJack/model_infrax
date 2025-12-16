#!/bin/bash

set -e  # 遇到错误时退出

mkdir -p model_infrax
cd model_infrax

echo "Downloading model_infrax..."
curl -LsSfO "https://github.com/LingoJack/model_infrax/raw/refs/heads/main/release/model_infrax.tar.gz"
tar -xzf model_infrax.tar.gz
rm -rf model_infrax.tar.gz

echo "Installing model_infrax..."
cd ..
mv model_infrax /Applications/
echo "model_infrax installed"

chmod +x /Applications/model_infrax/jen
chmod +x /Applications/model_infrax/jcode

echo 'export PATH="/Applications/model_infrax:$PATH"' >> ~/.zshrc
source ~/.zshrc

