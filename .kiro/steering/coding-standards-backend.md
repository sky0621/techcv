# Backendに関するコーディング規約

## 採用言語
- Golangの最新バージョンを採用する

## 基本原則

- `gofmt`でフォーマットされたコードを書く
- `golangci-lint`を使用して静的解析を行う
- シンプルで読みやすいコードを優先する

## 命名規則

### パッケージ名
- 小文字のみを使用し、アンダースコアやキャメルケースは使わない
- 短く簡潔な名前を使う（例: `http`, `user`, `auth`）
- 複数形ではなく単数形を使う

### 変数・関数名
- キャメルケースを使用する
- エクスポートする識別子は大文字で始める
- エクスポートしない識別子は小文字で始める
- 略語は全て大文字または全て小文字にする（例: `userID`, `HTTPServer`）

### インターフェース名
- 単一メソッドのインターフェースは `-er` サフィックスを使う（例: `Reader`, `Writer`）
- 意味のある名前を付ける

## コード構造

### ファイル構成
```go
package name

// imports
import (
    // 標準ライブラリ
    "context"
    "fmt"
    
    // 外部パッケージ
    "github.com/external/package"
    
    // 内部パッケージ
    "project/internal/domain"
)

// 定数
const (
    MaxRetries = 3
)

// 変数
var (
    ErrNotFound = errors.New("not found")
)

// 型定義
type User struct {}

// 関数
func NewUser() *User {}
```

### エラーハンドリング
- エラーは無視せず、必ず処理する
- エラーメッセージは小文字で始め、句読点で終わらない
- カスタムエラーは `errors.New()` または `fmt.Errorf()` で作成
- エラーのラップには `fmt.Errorf()` と `%w` を使用
```go
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}
```

### コンテキスト
- 関数の最初の引数として `context.Context` を受け取る
- コンテキストは構造体に保存しない
```go
func (s *Service) CreateUser(ctx context.Context, user *User) error {
    // ...
}
```

## ベストプラクティス

### 構造体
- ゼロ値で使える構造体を設計する
- コンストラクタ関数は `New` または `NewXxx` という名前にする
- フィールドは論理的な順序で並べる（重要なものを先に）

### インターフェース
- 小さなインターフェースを定義する
- 使う側でインターフェースを定義する（consumer-driven）
- 不要なインターフェースは作らない

### 並行処理
- goroutineのリークに注意する
- チャネルは送信側でクローズする
- `sync.WaitGroup` や `errgroup` を使って goroutine を管理する
- 共有メモリではなく、チャネルで通信する

### テスト
- テストファイルは `_test.go` サフィックスを使う
- テーブル駆動テストを活用する
- テスト関数名は `TestXxx` の形式にする
- サブテストには `t.Run()` を使用する
```go
func TestCreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   *User
        wantErr bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### コメント
- エクスポートされる識別子には必ずドキュメントコメントを書く
- コメントは識別子名で始める
```go
// User represents a user in the system.
type User struct {}

// CreateUser creates a new user with the given parameters.
func CreateUser(ctx context.Context, email string) (*User, error) {}
```

## 避けるべきパターン

- `panic` の使用（プログラムの初期化時以外）
- グローバル変数の多用
- `init()` 関数の過度な使用
- 不必要な抽象化
- 長すぎる関数（目安: 50行以内）
- 深いネスト（目安: 3レベル以内）

## 依存性注入

- コンストラクタで依存関係を注入する
- インターフェースを使って疎結合にする
```go
type Service struct {
    repo UserRepository
}

func NewService(repo UserRepository) *Service {
    return &Service{repo: repo}
}
```
