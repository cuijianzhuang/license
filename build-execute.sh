#!/bin/bash

cd web
npm install
npm run build
cd ../

# VERSION - 获取最新的tag并去除开头的'v'
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "0.0.1")
# 如果版本号以'v'开头，则去除这个'v'
VERSION=$(echo $VERSION | sed 's/^v//')
echo "Using version: $VERSION"

# 生成随机哈希值 (8位字母数字组合)
HASH=$(openssl rand -hex 4)
echo "Using hash: $HASH"

# 创建 build 目录，如果不存在的话
mkdir -p build

# Linux amd64
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X license/sys.Version=${VERSION} -X license/sys.Hash=${HASH} -X license/sys.Arch=linux/amd64" -o build/license-linux-amd64 .
# Linux arm64
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X license/sys.Version=${VERSION} -X license/sys.Hash=${HASH} -X license/sys.Arch=linux/arm64" -o build/license-linux-arm64 .
# macOS amd64
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X license/sys.Version=${VERSION} -X license/sys.Hash=${HASH} -X license/sys.Arch=darwin/amd64" -o build/license-darwin-amd64 .
# macOS arm64
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X license/sys.Version=${VERSION} -X license/sys.Hash=${HASH} -X license/sys.Arch=darwin/arm64" -o build/license-darwin-arm64 .
# Windows amd64
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X license/sys.Version=${VERSION} -X license/sys.Hash=${HASH} -X license/sys.Arch=windows/amd64" -o build/license-windows-amd64.exe .
# Windows arm64
GOOS=windows GOARCH=arm64 go build -ldflags="-s -w -X license/sys.Version=${VERSION} -X license/sys.Hash=${HASH} -X license/sys.Arch=windows/arm64" -o build/license-windows-arm64.exe .
