# Go Backend Bootstrap チェックリスト

## 対象

`services/<service>/backend` のような backend ディレクトリを初期構築する際にこのチェックリストを使う。

## ファイルセット

以下を作成または更新する:
- `go.mod`
- `.env.sample`
- `.gitignore`
- `Makefile`
- `README.md`
- `cmd/api/main.go`
- `internal/domain/health.go`
- `internal/usecase/health/service.go`
- `internal/interface/http/handler/health_handler.go`
- `internal/interface/http/server/server.go`
- `internal/infrastructure/config/config.go`
- `internal/infrastructure/clock/system_clock.go`
- `internal/infrastructure/sqlite/storage.go`

## 契約仕様

`/health` は次を含む JSON を返す:
- `"service"`: 固定のサービス名
- `"status"`: `"ok"`
- `"time"`: UTC の RFC3339 文字列

## Make ターゲット

以下を提供する:
- `make run`
- `make test`
- `make fmt`
- `make vet`
- `make tidy`
- `make healthcheck`

`make healthcheck` は status が `ok` でない場合に非ゼロ終了する。

## 検証

実行:

```bash
make test
make vet
```

必要に応じて実行:

```bash
make run
make healthcheck
```

サンドボックス制約でポート bind が失敗する場合は、その制約を明示してコンパイル/テスト結果を報告する。
