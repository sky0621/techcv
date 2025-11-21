# techcv
Webエンジニアの技術的な履歴（Curriculum Vitae）をまとめるプロジェクトです。

## バックエンドツール

`services/manager/backend` ディレクトリで以下を実行します。

- `make generate` – `docs/openapi.yaml` をもとに Echo 互換のハンドラや型を再生成します。
- `VERIFICATION_URL_BASE` – 任意。デフォルトは `http://localhost:5173/auth/verify` で、登録メールに記載する確認リンクの生成に使われます。

## Makefileの使い方

リポジトリルートの `Makefile` はディスパッチャとして動作し、各サービスとレイヤーの `Makefile` に処理を委譲します。`make <service>-<layer>-<goal>` の形式でターゲットを呼び出してください。

- `<service>`: `manager` / `publisher` / `administrator`
- `<layer>`: `backend` (略称 `be`) または `frontend` (略称 `fe`)
- `<goal>`: 各レイヤー配下の `Makefile` に定義されているターゲット名

例:

- `make manager-backend-test` – `services/manager/backend` 配下で `make test` を実行します。
- `make manager-backend-build` – manager バックエンドのバイナリをビルドします。
- `make manager-frontend-dev` – manager フロントエンドの開発用ターゲット（定義されていれば）を実行します。

OpenAPI 用のターゲットは `make openapi-<target>` で引き続き利用可能です（例: `make openapi-bundle-openapi` は `services/manager/openapi` で実行されます）。このエントリポイントを通してチーム全員が同じツールチェーンを利用してください。
