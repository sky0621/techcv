# テスト計画と運用

## テスト戦略
- **ユニットテスト**: 値オブジェクト（`GoogleID`）、ユーザーアグリゲート、`StartGoogleLogin` / `CompleteGoogleCallback` ユースケースの正常系・異常系を table-driven で検証する。
- **インテグレーションテスト**: MySQL リポジトリは `sqlmock` またはテスト用 DB を用い、Redis セッションリポジトリはコンテナ環境で検証する。Echo ハンドラは `httptest` でリクエスト/レスポンスを確認する。
- **契約テスト**: OpenAPI を redocly などで検証し、破壊的変更を防ぐ。生成コードの再生成が CI で落ちないよう `make generate-openapi` を追加する。
- **フロントエンドテスト**: React Testing Library + Vitest でコンポーネント・ページをテストし、MSW で API レスポンスをモックする。state 不一致や Google 拒否時の UI 表示を確認する。
- **E2E テスト（任意）**: テスト用 Google OAuth クレデンシャルもしくは擬似 OAuth サーバーを用いて、リダイレクトからダッシュボード表示までを確認する。

## ロールアウトと運用
- `docker-compose.yml` に Redis を追加する場合は、ローカル起動手順とクリーンアップ方法を README に追記する。
- バックエンドとフロントエンドのリダイレクト URI を揃え（例: `http://localhost:8080/techcv/api/v1/auth/google/callback` と `http://localhost:5173/auth/callback`）、Google Cloud Console に登録する。
- OAuth エラー率、トークン交換失敗、外部 API エラーを構造化ログに記録し、既存の監視にメトリクスを追加する。
- 機能フラグや UI の段階的公開が必要であれば、Google ログインボタン表示を設定値で切り替えられるようにする。
- インフラチームと協議し、`SameSite=None`、`Secure` 属性、クッキードメイン、シークレット管理を本番マニフェストに反映する。

## 未確定事項
- コールバックエンドポイントが JSON を返すか 302 リダイレクトを行うかを最終決定する必要がある（フロント実装との整合性が前提）。
- OAuth state の永続化方式（MySQL テーブル vs Redis）を確定し、将来のスケール要件に合わせたストアを採用する。
- 監査ログやアカウントリンク解除など、追加要件が存在する場合は別途タスクとして切り出す。
