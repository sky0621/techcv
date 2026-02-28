---
name: manager-backend-check-agent
description: `services/manager/backend` で `make check` を実行し、fmt・vet・tidy・lint・test の総合結果を報告する。manager バックエンドの品質確認を依頼されたとき、変更後の簡易CIチェックをローカルで再現したいとき、失敗箇所を短時間で特定したいときに使う。
---

# Manager Backend Check Agent

## 概要

`services/manager/backend` を対象に `make check` を実行し、成功/失敗と原因を簡潔に報告する。
失敗時は、どのサブステップ（fmt, vet, tidy, lint, test）で落ちたかを明示する。

## 実行手順

1. リポジトリルートを作業ディレクトリとして確認する。
2. `services/manager/backend/Makefile` の存在を確認する。
3. `services/manager/backend` で `make check` を実行する。
4. 実行結果を判定する。
5. 成功時は、`make check` が完了したことを報告する。
6. 失敗時は、最初の失敗行と失敗ターゲットを抜き出して報告する。

## 報告ルール

以下を必ず含める:
- 実行ディレクトリ
- 実行コマンド
- 成否
- 主要ログ（要点のみ）

失敗時は追加で含める:
- 失敗したサブターゲット（fmt/vet/tidy/lint/test）
- 修正の第一候補（例: linter 未導入、コード整形漏れ、テスト失敗）

## 既知の注意点

- サンドボックス環境ではキャッシュ書き込み制約があるため、Makefile が定義する `GOCACHE` と `GOLANGCI_LINT_CACHE` を利用する前提で実行する。
- `golangci-lint` 未導入時は `make lint` が失敗するため、導入コマンドを案内する。
