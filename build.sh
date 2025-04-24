#!/bin/bash

set -e

# 从 .env 文件中导入环境变量
if [ -f ".env" ]; then
    export $(cat .env | sed 's/#.*//g' | xargs)
else
    echo ".env file not found"
    exit 1
fi

# 使用环境变量中的用户名和密码尝试登录Docker Hub
docker login -u="${HUB_USER}" -p="${HUB_PASS}"
status=$?

# 检查登录命令的退出状态
if [ $status -ne 0 ]; then
    echo "Docker login failed, exiting..."
    exit $status
else
    echo "Docker login successful."
fi

# 创建并使用一个新的 Buildx 构建器实例，如果已存在则使用现有的
BUILDER_NAME=multi-platform-build
docker buildx create --name $BUILDER_NAME --use || true
docker buildx use $BUILDER_NAME
docker buildx inspect --bootstrap

# VERSION - 获取最新的tag并去除开头的'v'
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "0.0.1")
# 如果版本号以'v'开头，则去除这个'v'
VERSION=$(echo $VERSION | sed 's/^v//')
echo "Using version: $VERSION"

# 生成随机哈希值 (8位字母数字组合)
HASH=$(cat /dev/urandom | tr -dc 'a-f0-9' | fold -w 8 | head -n 1)
echo "Using commit hash: $HASH"

# 使用 Docker Buildx 构建镜像，同时标记为 latest 和 VERSION，支持多架构
docker buildx build \
  --no-cache \
  --platform linux/amd64,linux/arm64 \
  --build-arg VERSION=$VERSION \
  --build-arg HASH=$HASH \
  -t ${HUB_USER}/${HUB_REPO}:$VERSION \
  -t ${HUB_USER}/${HUB_REPO}:latest . \
  --push \
  --progress=plain

# 登出 Docker Hub
docker logout