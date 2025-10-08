# Requirements Document

## Introduction

WebエンジニアのCV（履歴書）を管理するシステムです。一般ユーザーがGoogle Accountで登録し、自身のCV情報を管理できます。CV情報は項目ごとに公開可否を設定でき、Web公開、Excel出力、PDF出力が可能です。

## Requirements

### Requirement 1: Google Accountでのユーザー登録

**User Story:** 一般ユーザーとして、自身のGoogle Accountでシステムの利用登録を行いたい。これにより、簡単かつ安全にシステムを利用開始できる。

#### Acceptance Criteria

1. WHEN ユーザーがシステムにアクセスする THEN システム SHALL Google OAuth 2.0認証画面を表示する
2. WHEN ユーザーがGoogle Accountで認証する THEN システム SHALL Googleアカウント情報（メールアドレス、名前）を取得する
3. WHEN 初回ログインユーザーの場合 THEN システム SHALL 新規ユーザーとしてデータベースに登録する
4. WHEN ユーザー登録が完了する THEN システム SHALL ユーザーをダッシュボード画面にリダイレクトする
5. WHEN 既存ユーザーがログインする THEN システム SHALL ユーザー情報を認証してダッシュボード画面にリダイレクトする
6. WHEN 認証に失敗する THEN システム SHALL エラーメッセージを表示して再認証を促す

### Requirement 2: CV情報の登録

**User Story:** 一般ユーザーとして、自身のCV情報を登録したい。これにより、自分の経歴やスキルをシステムで管理できる。

#### Acceptance Criteria

1. WHEN ユーザーがCV登録画面にアクセスする THEN システム SHALL CV情報入力フォームを表示する
2. WHEN CV入力フォームを表示する THEN システム SHALL 以下の項目を含む
   - 基本情報（氏名、メールアドレス、電話番号、住所、生年月日）
   - 職務経歴（会社名、在籍期間、役職、業務内容、使用技術スタック）
   - スキル情報（技術名、習熟度、経験年数）
   - 学歴（学校名、学部・学科、入学年、卒業年）
   - 資格・認定（資格名、取得年月日、認定機関）
   - プロジェクト実績（プロジェクト名、期間、役割、概要、使用技術）
3. WHEN ユーザーがCV情報を入力して保存ボタンをクリックする THEN システム SHALL 入力内容をバリデーションする
4. WHEN バリデーションが成功する THEN システム SHALL CV情報をデータベースに保存して成功メッセージを表示する
5. WHEN 必須項目が未入力の場合 THEN システム SHALL エラーメッセージを表示して保存を拒否する
6. WHEN 入力形式が不正な場合 THEN システム SHALL 該当項目にエラーメッセージを表示する
7. WHEN ユーザーが既存のCV情報を持つ場合 THEN システム SHALL 既存情報を表示して編集を許可する

### Requirement 3: CV項目ごとの公開可否設定

**User Story:** 一般ユーザーとして、自身のCVの項目ごとに公開可否を設定したい。これにより、状況に応じて適切な情報のみを公開できる。

#### Acceptance Criteria

1. WHEN ユーザーがCV編集画面にアクセスする THEN システム SHALL 各項目に公開/非公開の切り替えボタンを表示する
2. WHEN ユーザーが項目の公開設定を変更する THEN システム SHALL 変更内容を即座にデータベースに保存する
3. WHEN ユーザーが基本情報の特定項目を非公開に設定する THEN システム SHALL その項目の公開フラグをfalseに更新する
4. WHEN ユーザーが職務経歴の特定エントリを非公開に設定する THEN システム SHALL そのエントリの公開フラグをfalseに更新する
5. WHEN ユーザーがスキル情報の特定項目を非公開に設定する THEN システム SHALL その項目の公開フラグをfalseに更新する
6. WHEN ユーザーが学歴・資格・プロジェクト実績の特定項目を非公開に設定する THEN システム SHALL その項目の公開フラグをfalseに更新する
7. WHEN 公開設定の変更が完了する THEN システム SHALL 確認メッセージを表示する

### Requirement 4: CVのWeb公開

**User Story:** 一般ユーザーとして、自身のCVをWebに公開したい。これにより、他者がブラウザから自分のCVを閲覧できる。

#### Acceptance Criteria

1. WHEN ユーザーがCV公開設定を有効にする THEN システム SHALL 一意の公開URLを生成する
2. WHEN 公開URLが生成される THEN システム SHALL ユーザーにURLを表示してコピー機能を提供する
3. WHEN 第三者が公開URLにアクセスする THEN システム SHALL 公開設定された項目のみを含むCV情報を表示する
4. WHEN CV公開画面を表示する THEN システム SHALL レスポンシブデザインで見やすいレイアウトを提供する
5. WHEN 非公開に設定された項目がある THEN システム SHALL その項目を画面に表示しない
6. WHEN ユーザーがCV公開設定を無効にする THEN システム SHALL 公開URLへのアクセスを無効化する
7. WHEN 無効化されたURLにアクセスする THEN システム SHALL 「このCVは公開されていません」というメッセージを表示する

### Requirement 5: CVのExcel出力

**User Story:** 一般ユーザーとして、自身のCVをExcelファイルに出力したい。これにより、オフラインでの閲覧や編集、印刷が容易になる。

#### Acceptance Criteria

1. WHEN ユーザーがCV詳細画面でExcel出力ボタンをクリックする THEN システム SHALL 公開設定された項目のみを含むExcelファイルを生成する
2. WHEN Excelファイルを生成する THEN システム SHALL 読みやすいフォーマットで情報を整形する
3. WHEN Excelファイルを生成する THEN システム SHALL 各セクション（基本情報、職務経歴、スキル、学歴、資格、プロジェクト実績）を適切なシートまたはセクションに分割する
4. WHEN 非公開に設定された項目がある THEN システム SHALL その項目をExcelファイルに含めない
5. WHEN Excelファイルの生成が完了する THEN システム SHALL ファイルをダウンロードする
6. WHEN Excelファイル名を設定する THEN システム SHALL 「CV_{ユーザー名}_{日付}.xlsx」の形式で命名する

### Requirement 6: CVのPDF出力

**User Story:** 一般ユーザーとして、自身のCVをPDFファイルに出力したい。これにより、プロフェッショナルな形式で履歴書を共有できる。

#### Acceptance Criteria

1. WHEN ユーザーがCV詳細画面でPDF出力ボタンをクリックする THEN システム SHALL 公開設定された項目のみを含むPDFファイルを生成する
2. WHEN PDFファイルを生成する THEN システム SHALL プロフェッショナルで読みやすいレイアウトで情報を整形する
3. WHEN PDFファイルを生成する THEN システム SHALL 各セクション（基本情報、職務経歴、スキル、学歴、資格、プロジェクト実績）を適切に配置する
4. WHEN 非公開に設定された項目がある THEN システム SHALL その項目をPDFファイルに含めない
5. WHEN PDFファイルの生成が完了する THEN システム SHALL ファイルをダウンロードする
6. WHEN PDFファイル名を設定する THEN システム SHALL 「CV_{ユーザー名}_{日付}.pdf」の形式で命名する
7. WHEN PDFファイルを生成する THEN システム SHALL 印刷に適したA4サイズのページレイアウトを使用する
