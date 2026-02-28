# services-manager-backend
CV登録・編集機能を提供するサービスのバックエンド（Golang）のソースを管理するディレクトリです。

## アーキテクチャ
ヘキサゴナル・アーキテクチャを採用しています。

## 採用言語
Golang（ローカル確認: `go1.25.6`）

## データ永続化
SQLite3（開発用に `SQLITE_PATH` のファイル作成まで実装）

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
go test ./...
make run
```

## 開発コマンド
```bash
make run   # API起動
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
