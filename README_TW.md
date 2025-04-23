# 許可證管理服務

一個基於Go語言的軟體許可證驗證和認證管理服務。

![GitHub commit activity](https://img.shields.io/github/commit-activity/t/nannanStrawberry314/license/feture-go-v2?color=blue)
![GitHub forks](https://img.shields.io/github/forks/nannanStrawberry314/license?style=flat&color=brightgreen)
![GitHub stars](https://img.shields.io/github/stars/nannanStrawberry314/license?color=orange)
![GitHub pull requests](https://img.shields.io/github/issues-pr/nannanStrawberry314/license?color=red)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/nannanStrawberry314/license/.github/workflows/release.yml?label=build)
![GitHub Release](https://img.shields.io/github/v/release/nannanStrawberry314/license?color=brightgreen)
![Docker Pulls](https://img.shields.io/docker/pulls/raspberrycheese/license?color=blueviolet)

[简体中文](README.md) | [繁體中文](README_TW.md) | [Русский](README_RU.md) | [English](README_EN.md) | [日本語](README_JP.md) | [한국어](README_KR.md)

## 功能特點

- 各類軟體產品的許可證生成與驗證
- 支援JetBrains產品、GitLab、FinalShell、MobaXterm和JRebel等軟體
- 基於Gin框架構建的RESTful API介面
- 使用cron實現定時任務
- 通過GORM支援資料庫儲存(MySQL/SQLite)
- 採用RSA演算法的安全加密
- 多語言支援（簡體中文、繁體中文、俄語、英語、日語、韓語）

## 系統要求

- Go 1.24或更高版本
- MySQL資料庫(開發環境可使用SQLite)
- Docker(可選，用於容器化部署)

## 安裝方式

### 方式一：直接安裝

1. 克隆倉庫
   ```
   git clone https://github.com/nannanStrawberry314/license.git
   cd license
   ```

2. 安裝依賴
   ```
   go mod download
   ```

3. 配置環境變數(複製.env.example到.env並根據需要修改)

4. 構建並運行
   ```
   go build -o license-server
   ./license-server
   ```

### 方式二：Docker部署

1. 構建Docker映像
   ```
   docker build -t license-server .
   ```

2. 使用docker-compose運行
   ```
   docker-compose up -d
   ```

## 配置說明

配置通過環境變數和`.env`文件進行管理：

- `HTTP_HOST`：伺服器綁定的主機地址
- `HTTP_PORT`：監聽的連接埠
- `DB_TYPE`：資料庫類型(mysql或sqlite)
- `DB_DSN`：資料庫連接字串

## GitHub Actions配置

本專案包含GitHub Actions工作流，用於自動構建、測試和發布。要使用這些工作流，您需要配置以下倉庫密鑰：

- `HUB_USER`：您的Docker Hub用戶名
- `HUB_PASS`：您的Docker Hub密碼
- `HUB_REPO`：Docker Hub倉庫名稱
- `PUBLIC_REPO_TOKEN`：（可選）具有對公共倉庫寫入權限的個人訪問令牌
- `PUBLIC_REPO`：（可選）公共倉庫路徑，格式為`用戶名/倉庫名`

工作流包括：
- 每次推送時進行構建和測試
- 在標籤推送或手動觸發時構建和發布Docker映像
- 創建GitHub發布版本
- 將發布同步到公共倉庫（僅在手動觸發時執行，需要配置`PUBLIC_REPO_TOKEN`和`PUBLIC_REPO`）

## API介面

該服務提供多個許可證管理相關的API介面：

- `POST /v1/generate`：生成新的許可證
- `POST /v1/validate`：驗證現有許可證
- `GET /v1/status`：檢查服務狀態

詳細使用說明請參考API文檔。

## 開發指南

如需貢獻程式碼：

1. Fork本倉庫
2. 創建功能分支
3. 提交您的更改
4. 發起拉取請求

## 免責聲明

### 專案目的

本專案被開發和分享僅用於教育目的，旨在提供技術研究和學術學習的機會。專案內容涉及的技術和方法僅供學習和研究使用。

### 使用限制

專案作者發布此專案的目的不是鼓勵軟體盜版或任何形式的非法活動，而是為了促進知識的共享和技術的進步。**嚴禁將本專案用於破解、啟動或以任何方式非法使用軟體**。用戶在使用本專案時，必須遵守所有適用的當地和國際法律法規。

### 商業用途禁止

本專案**嚴禁用於商業用途**或任何形式的盈利活動。專案的任何部分都不得在任何可能直接或間接產生經濟利益的場合中使用。

### 免責聲明

本專案以「按原樣」方式提供，不附帶任何明示或暗示的保證。專案作者對於使用本專案內容可能導致的任何形式的損害或法律後果不承擔任何責任。使用本專案表示您理解並同意這些條件，並且您將自行承擔使用本專案的所有風險。

## Star歷史

[![Star History Chart](https://api.star-history.com/svg?repos=nannanStrawberry314/license&type=Date)](https://www.star-history.com/#nannanStrawberry314/license&Date) 