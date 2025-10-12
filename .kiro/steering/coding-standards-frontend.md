# Frontendに関するコーディング規約

## 採用技術スタック
- React 18以上（最新バージョン）
- TypeScript（厳格な型チェック）
- Vite（ビルドツール）
- TailwindCSS（スタイリング）

## 基本原則

- ESLintとPrettierでコードを統一
- TypeScriptの厳格モードを使用
- 関数コンポーネントとHooksを使用（クラスコンポーネントは使わない）
- 宣言的で読みやすいコードを優先
- 再利用可能なコンポーネントを設計

## 命名規則

### ファイル名
- コンポーネントファイルは PascalCase（例: `UserProfile.tsx`, `LoginForm.tsx`）
- Hooksファイルは camelCase で `use` プレフィックス（例: `useAuth.ts`, `useUserData.ts`）
- ユーティリティファイルは camelCase（例: `formatDate.ts`, `validation.ts`）
- 定数ファイルは camelCase（例: `constants.ts`, `apiEndpoints.ts`）

### 変数・関数名
- camelCase を使用（例: `userName`, `handleClick`）
- boolean値は `is`, `has`, `should` などのプレフィックス（例: `isLoading`, `hasError`）
- イベントハンドラーは `handle` プレフィックス（例: `handleSubmit`, `handleChange`）
- 定数は UPPER_SNAKE_CASE（例: `MAX_RETRY_COUNT`, `API_BASE_URL`）

### コンポーネント名
- PascalCase を使用
- 意味のある名前を付ける（例: `UserList`, `LoginButton`）
- ページコンポーネントは `Page` サフィックス（例: `UserListPage`, `LoginPage`）

### 型・インターフェース名
- PascalCase を使用
- Props型は `Props` サフィックス（例: `UserCardProps`, `ButtonProps`）
- インターフェースは `I` プレフィックスを使わない

## ディレクトリ構造

```
src/
├── components/          # 再利用可能なコンポーネント
│   ├── ui/             # 基本的なUIコンポーネント
│   │   ├── Button.tsx
│   │   └── Input.tsx
│   └── features/       # 機能別コンポーネント
│       └── user/
│           ├── UserCard.tsx
│           └── UserList.tsx
│
├── pages/              # ページコンポーネント
│   ├── LoginPage.tsx
│   └── UserListPage.tsx
│
├── hooks/              # カスタムHooks
│   ├── useAuth.ts
│   └── useUserData.ts
│
├── api/                # API通信
│   ├── client.ts       # APIクライアント設定
│   └── user.ts         # ユーザー関連API
│
├── types/              # 型定義
│   ├── user.ts
│   └── api.ts
│
├── utils/              # ユーティリティ関数
│   ├── formatDate.ts
│   └── validation.ts
│
├── constants/          # 定数
│   └── index.ts
│
└── App.tsx             # アプリケーションルート
```

## TypeScript

### 型定義
- `any` の使用は避ける
- 明示的な型注釈を付ける（推論できる場合は省略可）
- Props型は必ず定義する
```typescript
type UserCardProps = {
  user: User;
  onEdit: (id: string) => void;
  isLoading?: boolean;
};

export const UserCard = ({ user, onEdit, isLoading = false }: UserCardProps) => {
  // ...
};
```

### 型のエクスポート
- 型は `type` キーワードで定義（`interface` も可だが統一する）
- 再利用する型は `types/` ディレクトリで管理
```typescript
// types/user.ts
export type User = {
  id: string;
  email: string;
  name: string;
  createdAt: Date;
};

export type UserListResponse = {
  users: User[];
  total: number;
};
```

## コンポーネント設計

### 関数コンポーネント
- アロー関数で定義
- 名前付きエクスポートを使用
- Props型を明示
```typescript
type ButtonProps = {
  children: React.ReactNode;
  onClick: () => void;
  variant?: 'primary' | 'secondary';
  disabled?: boolean;
};

export const Button = ({ 
  children, 
  onClick, 
  variant = 'primary',
  disabled = false 
}: ButtonProps) => {
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={`btn btn-${variant}`}
    >
      {children}
    </button>
  );
};
```

### コンポーネントの分割
- 1ファイル1コンポーネントを原則とする
- 50行を超えたら分割を検討
- 再利用可能な部分は別コンポーネントに抽出

### Props
- デフォルト値はデストラクチャリングで設定
- オプショナルなPropsには `?` を使用
- children は `React.ReactNode` 型を使用

## Hooks

### カスタムHooks
- `use` プレフィックスで始める
- 1つの責務に集中
- 戻り値は配列またはオブジェクト
```typescript
export const useAuth = () => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // 認証状態の確認
  }, []);

  const login = async (email: string, password: string) => {
    // ログイン処理
  };

  const logout = () => {
    // ログアウト処理
  };

  return { user, isLoading, login, logout };
};
```

### Hooksの使用ルール
- コンポーネントのトップレベルでのみ呼び出す
- 条件分岐やループ内で呼び出さない
- カスタムHooks内でのみ他のHooksを呼び出す

### よく使うHooks
- `useState`: ローカルステート管理
- `useEffect`: 副作用の処理
- `useCallback`: 関数のメモ化
- `useMemo`: 値のメモ化
- `useRef`: DOM参照や値の保持

## 状態管理

### ローカルステート
- コンポーネント固有の状態は `useState` を使用
- 複雑な状態は `useReducer` を検討

### グローバルステート
- 認証状態などはContext APIを使用
- 大規模な場合はZustandやJotaiを検討
```typescript
// contexts/AuthContext.tsx
type AuthContextType = {
  user: User | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const auth = useAuth();
  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
};

export const useAuthContext = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuthContext must be used within AuthProvider');
  }
  return context;
};
```

## データフェッチング

### TanStack Query（React Query）
- サーバーステートの管理に使用
- キャッシュとリフェッチを自動管理
```typescript
// api/user.ts
export const fetchUsers = async (): Promise<User[]> => {
  const response = await fetch('/api/users');
  if (!response.ok) throw new Error('Failed to fetch users');
  return response.json();
};

// hooks/useUsers.ts
export const useUsers = () => {
  return useQuery({
    queryKey: ['users'],
    queryFn: fetchUsers,
  });
};

// components/UserList.tsx
export const UserList = () => {
  const { data: users, isLoading, error } = useUsers();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <ul>
      {users?.map(user => (
        <li key={user.id}>{user.name}</li>
      ))}
    </ul>
  );
};
```

### Mutation
```typescript
export const useCreateUser = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (user: CreateUserInput) => createUser(user),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
};
```

## スタイリング

### TailwindCSS
- ユーティリティクラスを使用
- カスタムクラスは最小限に
- 複雑なスタイルはコンポーネント化
```typescript
export const Card = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
      {children}
    </div>
  );
};
```

### 条件付きスタイル
- `clsx` または `classnames` を使用
```typescript
import { clsx } from 'clsx';

type ButtonProps = {
  variant: 'primary' | 'secondary';
  disabled?: boolean;
};

export const Button = ({ variant, disabled }: ButtonProps) => {
  return (
    <button
      className={clsx(
        'rounded px-4 py-2 font-medium',
        variant === 'primary' && 'bg-blue-600 text-white',
        variant === 'secondary' && 'bg-gray-200 text-gray-800',
        disabled && 'opacity-50 cursor-not-allowed'
      )}
    >
      Click me
    </button>
  );
};
```

## エラーハンドリング

### エラーバウンダリ
- 予期しないエラーをキャッチ
- ユーザーにわかりやすいメッセージを表示
```typescript
export class ErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean }
> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  render() {
    if (this.state.hasError) {
      return <div>エラーが発生しました</div>;
    }
    return this.props.children;
  }
}
```

### API エラー
- エラーレスポンスを適切に処理
- ユーザーにフィードバックを提供
```typescript
const { mutate, error } = useCreateUser();

const handleSubmit = (data: CreateUserInput) => {
  mutate(data, {
    onError: (error) => {
      if (error instanceof ApiError) {
        toast.error(error.message);
      } else {
        toast.error('予期しないエラーが発生しました');
      }
    },
  });
};
```

## パフォーマンス最適化

### メモ化
- 重い計算は `useMemo` でメモ化
- コールバック関数は `useCallback` でメモ化
```typescript
const expensiveValue = useMemo(() => {
  return computeExpensiveValue(data);
}, [data]);

const handleClick = useCallback(() => {
  doSomething(id);
}, [id]);
```

### コンポーネントの最適化
- `React.memo` で不要な再レンダリングを防ぐ
```typescript
export const UserCard = React.memo(({ user }: UserCardProps) => {
  return <div>{user.name}</div>;
});
```

### 遅延ロード
- ページコンポーネントは `React.lazy` で遅延ロード
```typescript
const UserListPage = lazy(() => import('./pages/UserListPage'));

<Suspense fallback={<div>Loading...</div>}>
  <UserListPage />
</Suspense>
```

## テスト

### テストファイル
- テストファイルは `*.test.tsx` または `*.spec.tsx`
- コンポーネントと同じディレクトリに配置

### Testing Library
- React Testing Libraryを使用
- ユーザーの視点でテスト
```typescript
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

describe('Button', () => {
  it('クリックイベントが発火する', async () => {
    const handleClick = jest.fn();
    render(<Button onClick={handleClick}>Click</Button>);
    
    await userEvent.click(screen.getByRole('button'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });
});
```

## アクセシビリティ

- セマンティックなHTMLを使用
- 適切なARIA属性を付与
- キーボード操作をサポート
- 十分なコントラスト比を確保
```typescript
<button
  onClick={handleClick}
  aria-label="ユーザーを削除"
  disabled={isLoading}
>
  削除
</button>
```

## ベストプラクティス

### DRY原則
- 重複するコードは共通化
- 再利用可能なコンポーネントを作成

### 単一責任の原則
- 1つのコンポーネントは1つの責務
- 複雑になったら分割

### Props Drilling の回避
- Context APIを活用
- 状態管理ライブラリを検討

### 型安全性
- `any` を避ける
- 厳格な型チェックを有効化
- ジェネリクスを活用

## 避けるべきパターン

- インラインスタイルの多用
- 巨大なコンポーネント（200行以上）
- Props Drillingの深いネスト
- useEffectの過度な使用
- グローバル変数の使用
- 直接的なDOM操作（refを使う場合を除く）

## コメント

- 複雑なロジックには説明を追加
- JSDocでコンポーネントの使い方を記述
```typescript
/**
 * ユーザー情報を表示するカードコンポーネント
 * 
 * @param user - 表示するユーザー情報
 * @param onEdit - 編集ボタンクリック時のコールバック
 * @param isLoading - ローディング状態
 */
export const UserCard = ({ user, onEdit, isLoading }: UserCardProps) => {
  // ...
};
```
