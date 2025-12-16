#!bin/bash
echo "installing..."
curl -LO https://github.com/LingoJack/model_infrax/raw/refs/heads/main/release/model_infrax.tar.gz
tar -xzf model_infrax.tar.gz -C /Applications/
rm -rf model_infrax.tar.gz
echo "install success"

echo "setting env..."
export PATH="/Applications/model_infrax:$PATH"
sudo echo "export PATH=\"/Applications/model_infrax:$PATH\"" >> ~/.zshrc
source ~/.zshrc

sudo chmod +x /Applications/model_infrax/jen
sudo chmox +x /Applications/model_infrax/jcode
echo "env set success"

echo "jen version: $(jen -v)"