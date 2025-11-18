go build -o jen main.go wire_gen.go
mkdir -p ~/dev/model_infrax/target/
mkdir -p ~/dev/model_infrax/target/assert/
mv jen ~/dev/model_infrax/target/
cp -r ~/dev/model_infrax/assert/prompt/ ~/dev/model_infrax/target/assert/
cp ~/dev/model_infrax/assert/application.yml ~/dev/model_infrax/target/
cp ~/dev/model_infrax/assert/schema.sql ~/dev/model_infrax/target/
cp ~/dev/model_infrax/assert/install.sh ~/dev/model_infrax/target/
cp ~/dev/model_infrax/assert/jenfile ~/dev/model_infrax/target/
mkdir -p ~/dev/model_infrax/target/output/
mkdir -p ~/dev/model_infrax/pack/

# 生成带时间戳的唯一包名，避免覆盖已有文件
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
PACKAGE_NAME="jen_${TIMESTAMP}.zip"
PACKAGE_PATH="$HOME/dev/model_infrax/pack/$PACKAGE_NAME"

# 切换到目标目录并创建压缩包
cd ~/dev/model_infrax/target && zip -r "$PACKAGE_PATH" .

# 输出打包结果信息
echo "📦 打包完成: $PACKAGE_NAME"
echo "📍 保存路径: $PACKAGE_PATH"
echo "📊 包大小: $(du -h "$PACKAGE_PATH" | cut -f1)"
echo "✅ 全新打包成功，不会覆盖已有文件"