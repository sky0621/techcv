---
name: go-backend-bootstrap
description: go.mod、Makefile、.env.sample、ヘキサゴナル構成、/health エンドポイント、make healthcheck を含む最小の Go バックエンド開発環境を作成・標準化する。services/*/backend に新規作成する場合、Go 雛形が欠落したディレクトリを再構築する場合、または GOCACHE をプロジェクト内に固定してサンドボックス安全にしたい場合に使う。
---

# Go Backend Bootstrap

## 概要

ヘルスチェックエンドポイントと `make healthcheck` を含む、実行可能かつテスト可能な Go バックエンド雛形を一度に整備する。

## ワークフロー

1. 作業前に、リポジトリ全体と対象ディレクトリの指示ファイルを確認する。
2. 対象の backend ディレクトリを調査し、既存のユーザー変更を保持する。
3. 基本ファイルを作成または更新する: `go.mod`, `.env.sample`, `.gitignore`, `Makefile`, `README.md`。
4. 最小のヘキサゴナル構成を作成する: `cmd/`, `internal/domain/`, `internal/usecase/`, `internal/interface/http/`, `internal/infrastructure/`。
5. service 名、status、UTC タイムスタンプを返す `/health` エンドポイントを追加する。
6. `Makefile` に `run`, `test`, `fmt`, `vet`, `tidy`, `healthcheck` を追加する。
7. サンドボックス安全に実行できるよう `GOCACHE` をプロジェクト内パスへ固定する。
8. `make test` と `make vet` を実行する。
9. 実行確認を試み、環境制約でポート bind ができない場合は制約を明示して報告する。

## 必須 Makefile ルール

定義:
- `GOCACHE ?= $(CURDIR)/.cache/go-build`
- `HEALTHCHECK_URL ?= http://127.0.0.1:$(APP_PORT)/health`

実装:
- `run`: `APP_ENV`, `APP_PORT`, `SQLITE_PATH`, `GOCACHE` を付与して API を起動する。
- `test`: `GOCACHE` を指定して `go test ./...` を実行する。
- `vet`: `GOCACHE` を指定して `go vet ./...` を実行する。
- `healthcheck`: ヘルスチェック URL を `curl -fsS` で叩き、`"status":"ok"` が無ければ失敗終了する。

`.cache/`, `tmp/`, `.env`, カバレッジ成果物は `.gitignore` に含める。

## 参照

ファイル単位の受け入れ条件と期待レイアウトは `references/checklist.md` を参照する。
