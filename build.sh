go build -o model_infrax main.go wire_gen.go
mkdir -p ~/dev/model_infrax/cmd/
mkdir -p ~/dev/model_infrax/cmd/assert/
mv model_infrax ~/dev/model_infrax/cmd/
cp -r ~/dev/model_infrax/assert/prompt/ ~/dev/model_infrax/cmd/assert/
cp ~/dev/model_infrax/assert/application.yml ~/dev/model_infrax/cmd/
cp ~/dev/model_infrax/assert/schema.sql ~/dev/model_infrax/cmd/
cp ~/dev/model_infrax/assert/install.sh ~/dev/model_infrax/cmd/
cp ~/dev/model_infrax/assert/jcode ~/dev/model_infrax/cmd/
mkdir -p ~/dev/model_infrax/cmd/output/
mkdir -p ~/dev/model_infrax/pack/

# 创建压缩包，如果文件存在则覆盖
PACKAGE_NAME="model_infrax.zip"
PACKAGE_PATH="$HOME/dev/model_infrax/pack/$PACKAGE_NAME"

# 切换到目标目录并创建压缩包
cd ~/dev/model_infrax/cmd && zip -r "$PACKAGE_PATH" .

# 输出打包结果信息
echo "📦 打包完成: $PACKAGE_NAME"
echo "📍 保存路径: $PACKAGE_PATH"
echo "📊 包大小: $(du -h "$PACKAGE_PATH" | cut -f1)"
echo "✅ 打包成功"