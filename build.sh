go build -o jen cmd/jen/main.go cmd/jen/wire_gen.go
mkdir -p ~/dev/model_infrax/target/
mkdir -p ~/dev/model_infrax/target/assets/
mv jen ~/dev/model_infrax/target/
cp -r ~/dev/model_infrax/assets/prompt/ ~/dev/model_infrax/target/assets/
cp ~/dev/model_infrax/assets/application.yml ~/dev/model_infrax/target/
cp ~/dev/model_infrax/assets/schema.sql ~/dev/model_infrax/target/
cp ~/dev/model_infrax/assets/install.sh ~/dev/model_infrax/target/
cp ~/dev/model_infrax/assets/jcode ~/dev/model_infrax/target/
mkdir -p ~/dev/model_infrax/target/output/
mkdir -p ~/dev/model_infrax/pack/

# 创建压缩包，如果文件存在则覆盖
PACKAGE_NAME="model_infrax.zip"
PACKAGE_PATH="$HOME/dev/model_infrax/pack/$PACKAGE_NAME"

# 切换到目标目录并创建压缩包
cd ~/dev/model_infrax/target && zip -r "$PACKAGE_PATH" .

# 输出打包结果信息
echo "📦 打包完成: $PACKAGE_NAME"
echo "📍 保存路径: $PACKAGE_PATH"
echo "📊 包大小: $(du -h "$PACKAGE_PATH" | cut -f1)"
echo "✅ 打包成功"