# services-manager-backend
CV登録・編集機能を提供するサービスのバックエンド（Golang）のソースを管理するディレクトリです。

## アーキテクチャ
クリーン・アーキテクチャを採用しています。

## 採用言語
Golang（ローカル確認: `go1.25.6`）

## データ永続化
SQLite3（`POST /profiles` は `SQLITE_PATH` のDBへ永続化）

## DB関連ライブラリ
- データベースアクセス: `sqlc`（生成コード）
- データベースマイグレーション: `sqldef`（`sqlite3def`）

## ディレクトリ構成（最小雛形）
```text
backend/
  cmd/api/                      # エントリポイント
  internal/
    domain/                     # ドメイン層
    usecase/                    # ユースケース層
    interface/http/             # HTTP公開層
    infrastructure/             # 設定・時計・SQLite準備
  .env.sample
  Makefile
  go.mod
```

## セットアップ
```bash
cd services/manager/backend
cp .env.sample .env
go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.24.0
go install github.com/sqldef/sqldef/cmd/sqlite3def@latest
make sqlc-generate
go test ./...
make run
```

## 開発コマンド
```bash
make run   # API起動
make db-migrate # sqldefでスキーマ反映
make sqlc-generate # sqlc生成コード更新
make healthcheck # /health の status=ok を確認（API起動後）
make test  # 単体テスト
make fmt   # gofmt
make vet   # go vet
make lint  # golangci-lint
make tidy  # go.mod/go.sum整理
make check # fmt, vet, tidy, lint, test を順に実行
```

`make lint` の実行には `golangci-lint` が必要です。未導入の場合:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8
```

## 動作確認
起動後、以下で疎通確認できます。
```bash
make healthcheck
```

`make run` は起動前に `make db-migrate` を実行します。

名前とニックネームの登録は以下で確認できます。
```bash
curl -X POST http://localhost:8080/profiles \
  -H "Content-Type: application/json" \
  -d '{"name":"Taro","nickname":"taro-dev"}'
```

登録データは SQLite の `profiles` テーブルに保存されます。
