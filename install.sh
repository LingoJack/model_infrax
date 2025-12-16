#!bin/bash
echo "installing..."
curl -LO https://github.com/LingoJack/model_infrax/blob/main/release/model_infrax.tar.gz
tar -xzf model_infrax.tar.gz -C /Applications/
rm -rf model_infrax.tar.gz
echo "install success"

echo "setting env..."
export PATH="/Applications/model_infrax:$PATH"
echo "export PATH=\"/Applications/model_infrax:$PATH\"" >> ~/.zshrc
source ~/.zshrc
echo "env set success"

echo "jen version: $(jen -v)"