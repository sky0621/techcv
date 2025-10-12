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

## ADR
- ID
  - 各種IDには数値ではなくUUID v7を採用する

- go.modのtoolディレクティブ
  - OpenAPIソース自動生成やDBマイグレーション等、開発用の各種ツールの管理はgo.modのtoolディレクティブを使用して行う

- ログフォーマット
  - CloudLoggingに沿ったログ構造とする
    - フォーマット
      - JSONで構造化されフォーマットとする
    - ログレベル
      - DEBUG/INFO/WARN/ERRORなど

- API設計
  - RESTを採用する
    - ただしリソースベースでは扱いが難しい、パフォーマンスに影響出る場合など、はREST原則を崩す事も認める
  - エンドポイント
    - `/techcv/api/v1` をベースにする

- OpenAPI仕様の採用バージョン
  - OpenAPI 3.0 を採用する

- DIライブラリ
  - 現時点では導入しない

- レイヤー別の単体テスト方法
  - infrastructure
    - DBに対してはモックを使わずテストコードを書く
  - domain/usecase
    - infrastructureレイヤーをモックにしてテストコードを書く
  - adapter
    - 必要に応じて下位レイヤーをモックにしてテストコードを書く

- OpenAPIファイルの分割管理
  - OpenAPIファイルを分割して管理する
  - 分割方法
    - root.yaml
    - components
      - paths.yaml
      - parameters.yaml
      - responses.yaml
      - schemas.yaml

- リクエストパラメーターのバリデーション
  - 型、桁レベル（や必須等）のAPIスキーマでチェック可能なものはOpenAPIのYamlの方に記載してチェックする
  - 上記以外のドメインに関するものはドメイン層でチェックする

- 設定
  - 環境変数に保存し環境ごとのグルーピングをしない
  - ただし、ローカル環境用だけは個別にファイルを用意して読み込む方式でもよい

- 時刻の扱い
  - プロジェクト全体として日付時刻はUTCで統一して扱い、表示や計算で他のタイムゾーンが必要になった際にUTCから変換して処理を行う
  - アプリケーションサーバー
    - サーバーのタイムゾーンはUTCを利用する
    - goで扱う時刻のタイムゾーンはUTCに統一する
  - DB
    - DBサーバーのタイムゾーンはUTCを利用す
  - サーバー外部との時刻の入出力(APIのI/Fの形式など)
    - ISO8601 拡張形式を使用する
    - タイムゾーンはUTCを使用する
    - 精度はミリ秒単位まで扱う

- slogでのログ出力
  - Context付きログ関数の使用を必須とする
  - カスタムハンドラーによる自動情報付与
    - context.Contextに含まれるJWT由来の情報 (ユーザー識別用のID等) を自動的にログに追加
  - request_idなどのリクエスト追跡情報も自動付与

- HTTPレスポンス構造
  - すべての成功レスポンス（2xx系ステータスコード）はエンドポイントによって個別に定義する
  - エラーレスポンス構造
```
{
  "requestId": "88374925",
  "code": "VALIDATION_ERROR",
  "details": [
    {
      "field": "email",
      "code": "INVALID_EMAIL_FORMAT"
    }
  ]
}
```
    - requestId: 必須 ... １リクエストをユニークに特定するためのランダムID
    - code: omitempty ... エラーコード（大文字のスネークケース）
    - details: omitempty ... 詳細エラー情報の配列（任意、主にバリデーションエラー用）


