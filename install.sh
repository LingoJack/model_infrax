#!/bin/bash

set -e  # 遇到错误时退出

mkdir -p model_infrax
cd model_infrax

curl -O "https://github.com/LingoJack/model_infrax/raw/refs/heads/main/release/model_infrax.tar.gz"
tar -xzf model_infrax.tar.gz
rm model_infrax.tar.gz

tree .
