# ライセンス管理サービス

Goベースのソフトウェアライセンス検証と認証管理サービス。

![GitHub commit activity](https://img.shields.io/github/commit-activity/t/nannanStrawberry314/license?color=blue)
![GitHub forks](https://img.shields.io/github/forks/nannanStrawberry314/license?style=flat&color=brightgreen)
![GitHub stars](https://img.shields.io/github/stars/nannanStrawberry314/license?color=orange)
![GitHub pull requests](https://img.shields.io/github/issues-pr/nannanStrawberry314/license?color=red)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/nannanStrawberry314/license/.github/workflows/release.yml?label=build)
![GitHub Release](https://img.shields.io/github/v/release/nannanStrawberry314/license?color=brightgreen)
![Docker Pulls](https://img.shields.io/docker/pulls/raspberrycheese/license?color=blueviolet)

[简体中文](README_CN.md) | [繁體中文](README_TW.md) | [Русский](README_RU.md) | [English](README.md) | [日本語](README_JP.md) | [한국어](README_KR.md)

## 機能

- 様々なソフトウェア製品のライセンス生成と検証
- JetBrains製品、GitLab、FinalShell、MobaXterm、JRebelなどのサポート
- GinフレームワークによるRESTful APIインターフェース
- cronによるスケジュールタスク
- GORM（MySQL/SQLiteサポート）によるデータベース保存
- RSAによる安全な暗号化
- 多言語サポート（簡体字中国語、繁体字中国語、ロシア語、英語、日本語、韓国語）

## 必要条件

- Go 1.24以上
- MySQLデータベース（開発環境ではSQLiteも可）
- Docker（オプション、コンテナ化デプロイメント用）

## インストール方法

### 方法1：直接インストール

1. リポジトリをクローン
   ```
   git clone https://github.com/nannanStrawberry314/license.git
   cd license
   ```

2. 依存関係をインストール
   ```
   go mod download
   ```

3. 環境変数を設定（.env.exampleを.envにコピーして必要に応じて変更）

4. ビルドと実行
   ```
   go build -o license-server
   ./license-server
   ```

### 方法2：Dockerデプロイ

1. Dockerイメージをビルド
   ```
   docker build -t license-server .
   ```

2. docker-composeで実行
   ```
   docker-compose up -d
   ```

## 設定

設定は環境変数と`.env`ファイルで管理されます：

- `HTTP_HOST`：サーバーのホストアドレス
- `HTTP_PORT`：リッスンするポート
- `DB_TYPE`：データベースタイプ（mysqlまたはsqlite）
- `DB_DSN`：データベース接続文字列

## GitHub Actions設定

このプロジェクトには、自動ビルド、テスト、リリース用のGitHub Actionsワークフローが含まれています。これらのワークフローを使用するには、以下のリポジトリシークレットを設定する必要があります：

- `HUB_USER`：Docker Hubのユーザー名
- `HUB_PASS`：Docker Hubのパスワード
- `HUB_REPO`：Docker Hubリポジトリ名
- `PUBLIC_REPO_TOKEN`：（オプション）公開リポジトリへの書き込み権限を持つパーソナルアクセストークン
- `PUBLIC_REPO`：（オプション）`ユーザー名/リポジトリ`形式の公開リポジトリパス

ワークフローには以下が含まれます：
- 各プッシュでのビルドとテスト
- タグプッシュまたは手動トリガーでのDockerイメージビルドと公開
- GitHubリリースの作成
- 公開リポジトリへのリリースの同期（手動トリガー時のみ、`PUBLIC_REPO_TOKEN`と`PUBLIC_REPO`が必要）

## APIエンドポイント

APIは様々なライセンス管理エンドポイントを提供します：

- `POST /v1/generate`：新しいライセンスを生成
- `POST /v1/validate`：既存のライセンスを検証
- `GET /v1/status`：サービスの状態を確認

詳細な使用方法についてはAPIドキュメントを参照してください。

## 開発

このプロジェクトに貢献するには：

1. リポジトリをフォーク
2. 機能ブランチを作成
3. 変更を加える
4. プルリクエストを提出

## 免責事項

### プロジェクトの目的

このプロジェクトは教育目的のみのために開発・共有されており、技術研究と学術学習の機会を提供することを目的としています。このプロジェクトで議論されている技術や方法は、学習・研究目的のみを意図しています。

### 使用制限

このプロジェクトの作者は、ソフトウェア海賊行為や違法行為を奨励するためではなく、知識共有と技術進歩を促進するためにこのプロジェクトを公開しています。**本プロジェクトを使用してソフトウェアを違法にクラック、アクティベート、使用することは厳禁です**。ユーザーは本プロジェクトを使用する際、適用されるすべての地域および国際法・規制を遵守する必要があります。

### 商用利用禁止

本プロジェクトは**商用利用または利益を生み出す活動に使用することを厳禁**します。直接または間接的に経済的利益を生じさせる可能性のある状況でプロジェクトのいかなる部分も使用してはなりません。

### 責任の免除

本プロジェクトは「現状のまま」提供され、明示または黙示を問わず、いかなる保証もありません。プロジェクトの作者は、本プロジェクトのコンテンツを使用することによって生じる可能性のあるいかなる形の損害や法的結果に対しても責任を負いません。本プロジェクトを使用することにより、あなたはこれらの条件を理解し同意するとともに、本プロジェクトの使用に関するすべてのリスクを自己負担することになります。

## スター履歴

[![Star History Chart](https://api.star-history.com/svg?repos=nannanStrawberry314/license&type=Date)](https://www.star-history.com/#nannanStrawberry314/license&Date) 