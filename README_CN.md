# 许可证管理服务

一个基于Go语言的软件许可证验证和认证管理服务。

![GitHub commit activity](https://img.shields.io/github/commit-activity/t/nannanStrawberry314/license?color=blue)
![GitHub forks](https://img.shields.io/github/forks/nannanStrawberry314/license?style=flat&color=brightgreen)
![GitHub stars](https://img.shields.io/github/stars/nannanStrawberry314/license?color=orange)
![GitHub pull requests](https://img.shields.io/github/issues-pr/nannanStrawberry314/license?color=red)
![Docker Pulls](https://img.shields.io/docker/pulls/raspberrycheese/license?color=blueviolet)

[English](README.md) | [中文](README_CN.md)

## 功能特点

- 各类软件产品的许可证生成与验证
- 支持JetBrains产品、GitLab、FinalShell、MobaXterm和JRebel等软件
- 基于Gin框架构建的RESTful API接口
- 使用cron实现定时任务
- 通过GORM支持数据库存储(MySQL/SQLite)
- 采用RSA算法的安全加密

## 系统要求

- Go 1.21或更高版本
- MySQL数据库(开发环境可使用SQLite)
- Docker(可选，用于容器化部署)

## 安装方式

### 方式一：直接安装

1. 克隆仓库
   ```
   git clone https://github.com/nannanStrawberry314/license.git
   cd license
   ```

2. 安装依赖
   ```
   go mod download
   ```

3. 配置环境变量(复制.env.example到.env并根据需要修改)

4. 构建并运行
   ```
   go build -o license-server
   ./license-server
   ```

### 方式二：Docker部署

1. 构建Docker镜像
   ```
   docker build -t license-server .
   ```

2. 使用docker-compose运行
   ```
   docker-compose up -d
   ```

## 配置说明

配置通过环境变量和`.env`文件进行管理：

- `HTTP_HOST`：服务器绑定的主机地址
- `HTTP_PORT`：监听的端口
- `DB_TYPE`：数据库类型(mysql或sqlite)
- `DB_DSN`：数据库连接字符串

## API接口

该服务提供多个许可证管理相关的API接口：

- `POST /v1/generate`：生成新的许可证
- `POST /v1/validate`：验证现有许可证
- `GET /v1/status`：检查服务状态

详细使用说明请参考API文档。

## 开发指南

如需贡献代码：

1. Fork本仓库
2. 创建功能分支
3. 提交您的更改
4. 发起拉取请求

## 免责声明

### 项目目的

本项目被开发和分享仅用于教育目的，旨在提供技术研究和学术学习的机会。项目内容涉及的技术和方法仅供学习和研究使用。

### 使用限制

项目作者发布此项目的目的不是鼓励软件盗版或任何形式的非法活动，而是为了促进知识的共享和技术的进步。**严禁将本项目用于破解、激活或以任何方式非法使用软件**。用户在使用本项目时，必须遵守所有适用的当地和国际法律法规。

### 商业用途禁止

本项目**严禁用于商业用途**或任何形式的盈利活动。项目的任何部分都不得在任何可能直接或间接产生经济利益的场合中使用。

### 免责声明

本项目以"按原样"方式提供，不附带任何明示或暗示的保证。项目作者对于使用本项目内容可能导致的任何形式的损害或法律后果不承担任何责任。使用本项目表示您理解并同意这些条件，并且您将自行承担使用本项目的所有风险。

## Star历史

[![Star History Chart](https://api.star-history.com/svg?repos=nannanStrawberry314/license&type=Date)](https://www.star-history.com/#nannanStrawberry314/license&Date) 