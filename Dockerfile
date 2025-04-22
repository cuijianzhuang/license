FROM node:22-alpine AS frontend-builder
WORKDIR /app
# 复制package.json和package-lock.json以利用缓存
COPY web/package*.json ./
# 安装依赖
RUN npm install
# 复制React前端源代码
COPY web/ ./
# 构建React应用
RUN npm run build

FROM golang:1.24 AS go-builder
WORKDIR /app
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
# 复制Go源代码
COPY . .
# 删除可能存在的web/build目录，避免嵌入旧文件
RUN rm -rf web/build
# 复制前端构建产物
COPY --from=frontend-builder /app/build/ ./web/build/
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=1
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-s -w -extldflags '-static'" -o license ./main.go

FROM alpine AS runner
WORKDIR /app
COPY --from=go-builder /app/license ./license
RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true \
RUN mkdir -p /data

RUN ["chmod", "+x", "/app/license"]
ENTRYPOINT ["/app/license"]

