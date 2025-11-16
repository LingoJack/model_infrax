go build -o jen main.go wire_gen.go
mkdir -p ~/dev/model_infrax/target/
mkdir -p ~/dev/model_infrax/target/assert/
mv jen ~/dev/model_infrax/target/
cp -r ~/dev/model_infrax/assert/ ~/dev/model_infrax/target/assert/
zip -r ~/dev/model_infrax/target/jen.zip ./target/jen ./target/assert