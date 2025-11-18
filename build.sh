go build -o jen main.go wire_gen.go
mkdir -p ~/dev/model_infrax/target/
mkdir -p ~/dev/model_infrax/target/assert/
mv jen ~/dev/model_infrax/target/
cp -r ~/dev/model_infrax/assert/prompt/ ~/dev/model_infrax/target/assert/
cp ~/dev/model_infrax/assert/application.yml ~/dev/model_infrax/target/
cp ~/dev/model_infrax/assert/schema.sql ~/dev/model_infrax/target/
cp ~/dev/model_infrax/assert/install.sh ~/dev/model_infrax/target/
mkdir -p ~/dev/model_infrax/target/output/
mkdir -p ~/dev/model_infrax/pack/
cd ~/dev/model_infrax/target && zip -r ~/dev/model_infrax/pack/jen.zip .