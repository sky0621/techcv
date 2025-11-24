# Backendコーディング時のルール

- openapi.yamlを直接修正しない
  - 修正するのはspec配下のYamlファイル
  - openapi.yamlはredoclyによって自動生成される対象

- sqlcによって生成されたファイルを直接修正しない

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
  - OpenAPI 3.0.3 を採用する

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

- ログ出力
  - 極力、slog を使う。
  - fmt.Println や fmt.Fprintf などは使わない。

- slogでのログ出力
  - Context付きログ関数の使用を必須とする
  - カスタムハンドラーによる自動情報付与
    - context.Contextに含まれるJWT由来の情報 (ユーザー識別用のID等) を自動的にログに追加
  - request_idなどのリクエスト追跡情報も自動付与

- HTTPレスポンス構造
  - すべての成功レスポンス（2xx系ステータスコード）はエンドポイントによって個別に定義する
  - エラーレスポンス構造は以下の定義とする
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


