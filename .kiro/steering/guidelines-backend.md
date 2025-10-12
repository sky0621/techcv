# Backendのプロジェクトガイドライン

## プロジェクト概要
WebエンジニアのCV（履歴書＋職務経歴書）の管理やシステム全体に関わる設定等を行うためのWebAPIを提供する

## 技術スタック
- プログラミング言語
  - Golang 1.25
- Webフレームワーク
  - Echo
- API仕様
  - OpenAPI 3.0.3
- OpenAPIライブラリ
  - kin-openapi (github.com/getkin/kin-openapi)
- OpenAPIコード生成
  - oapi-codegen (github.com/oapi-codegen/oapi-codegen)
- データベースアクセス
  - sqlc
- データベースマイグレーション
  - sqldef
- データベース
  - MySQL 8.0+
- 認証
  - Google OAuth 2.0 (golang.org/x/oauth2)
- 環境変数管理
  - envconfig (github.com/vrischmann/envconfig)
- ローカル環境変数
  - godotenv (github.com/joho/godotenv)
- ログ
  - slog (Go標準ライブラリ)
- タスクランナー
  - Makefile
- Linter
  - golangci-lint

