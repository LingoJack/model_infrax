go build -o jen main.go wire_gen.go
mkdir -p ~/dev/model_infrax/target/
mkdir -p ~/dev/model_infrax/target/assert/
mv jen ~/dev/model_infrax/target/
cp -r ~/dev/model_infrax/assert/ ~/dev/model_infrax/target/assert/
mv ~/dev/model_infrax/target/assert/application.yml ~/dev/model_infrax/target/
mv ~/dev/model_infrax/target/assert/schema.sql ~/dev/model_infrax/target/
mv ~/dev/model_infrax/target/assert/install.sh ~/dev/model_infrax/target/
mkdir -p ~/dev/model_infrax/target/output/
cd ~/dev/model_infrax/target && zip -r ~/dev/model_infrax/pack/jen.zip .