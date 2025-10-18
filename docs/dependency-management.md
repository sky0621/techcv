# 依存関係管理

このリポジトリでは Renovate を利用して各サービスの依存関係を更新します。設定ファイルはリポジトリ直下の `renovate.json` で管理しており、以下の領域を対象にしています。

- `services/manager/backend` 配下の Go モジュール
- `services/manager/frontend` と `services/manager/openapi` の npm 依存関係

Renovate が作成する Pull Request には `dependencies` ラベルが自動付与されます。また依存関係ダッシュボード Issue が自動生成され、手動での再チェックや保留中の更新状況の確認に利用できます。

現在の設定では、新しいリリースが出てから 7 日間は Pull Request を作成しません。安定性を優先したい場合に便利ですが、より早く更新したい場合は `minimumReleaseAge` を調整してください。

更新頻度の変更や新しいパッケージグループの追加が必要な場合は `renovate.json` を編集し、デフォルトブランチへ反映してください。詳細なオプションについては [Renovate のドキュメント](https://docs.renovatebot.com/) を参照してください。
