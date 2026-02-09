# CLAUDE.md

## プロジェクト概要

mdprev — ローカルMarkdownプレビューツール。Go バックエンド + React フロントエンド。

仕様書: `docs/MVP_SPEC.md`

## リポジトリ構成

```
cmd/mdprev/           # Go エントリポイント
internal/
  domain/             # ドメイン層（エンティティ、値オブジェクト）
  usecase/            # ユースケース層
  interface/
    handler/          # HTTPハンドラ
    middleware/       # ミドルウェア（セキュリティ等）
  infrastructure/     # インフラ層（ファイルシステム操作）
web/                  # React フロントエンド（Vite + pnpm）
web/dist/             # ビルド成果物（.gitignore 対象）
docs/                 # 仕様書
```

## アーキテクチャ

### バックエンド（Go）

- **Clean Architecture + DDD**
- 依存方向: handler → usecase → domain ← infrastructure
- domain 層は外部依存を持たない
- infrastructure 層が domain のインターフェースを実装する

```
handler (interface層)
  ↓ 呼び出し
usecase (ユースケース層)
  ↓ 依存
domain (ドメイン層: エンティティ, リポジトリインターフェース)
  ↑ 実装
infrastructure (インフラ層: ファイルシステム)
```

### フロントエンド（React）

- **Vertical Slice Architecture**
- 機能単位でディレクトリを切る（共有部分は `shared/` に配置）

```
web/src/
  features/
    tree/             # ツリー機能（コンポーネント, hooks, API呼び出し）
    preview/          # プレビュー機能
    pathbar/          # パスバー機能
  shared/
    components/       # 共有UIコンポーネント
    hooks/            # 共有hooks
    lib/              # ユーティリティ
    types/            # 共有型定義
  App.tsx
  main.tsx
```

## 開発環境

### バージョン管理: mise

```toml
# .mise.toml
[tools]
go = "1.24"
node = "22"
pnpm = "10"
```

- `mise install` で Go, Node.js, pnpm をセットアップ

### バックエンド

- **言語**: Go 1.24
- **ルーター**: 標準ライブラリ (`net/http`)

### フロントエンド

- **言語**: TypeScript (strict mode)
- **フレームワーク**: React
- **ビルド**: Vite
- **パッケージマネージャ**: pnpm
- **CSS**: Tailwind CSS

## Linter / Formatter

### フロントエンド: Biome

```jsonc
// web/biome.json
{
  "formatter": {
    "indentStyle": "space",
    "indentWidth": 2
  },
  "linter": {
    "enabled": true
  }
}
```

- **フォーマット**: `pnpm biome format --write .`
- **リント**: `pnpm biome lint .`
- **両方**: `pnpm biome check --write .`

### バックエンド: gofmt / go vet

- **フォーマット**: `gofmt -w .`
- **静的解析**: `go vet ./...`

## よく使うコマンド

```bash
# --- セットアップ ---
mise install                      # Go, Node.js, pnpm インストール
cd web && pnpm install             # フロントエンド依存インストール

# --- フロントエンド ---
cd web && pnpm dev                 # Vite 開発サーバ起動
cd web && pnpm build               # web/dist/ にビルド
cd web && pnpm biome check --write . # lint + format

# --- バックエンド ---
go run ./cmd/mdprev                # サーバ起動
go test ./...                      # テスト実行
gofmt -w .                         # フォーマット
go vet ./...                       # 静的解析

# --- ビルド ---
cd web && pnpm build && cd .. && go build -o mdprev ./cmd/mdprev
```

## コーディング規約

### Go

- 標準の Go スタイルに従う（gofmt 準拠）
- エラーは呼び出し元に返す（`fmt.Errorf("xxx: %w", err)`）
- domain 層に外部パッケージの依存を入れない

### TypeScript / React

- Biome のルールに従う
- コンポーネントは関数コンポーネント + hooks
- 型は `interface` 優先（`type` は union など必要な場合のみ）
- feature 内で完結する処理は feature ディレクトリ内に閉じる
