# LLinter - Go Import Dependency Linter

LLinterはGoのimport文をチェックし、許可されていない依存関係を検出するLinterです。

## 特徴

- ファイルパスパターンに基づいたルール適用
- 禁止（deny）と許可（allow）のimportパターン設定
- YAMLベースの設定ファイル
- golang.org/x/tools/go/analysisフレームワークを使用

## インストール

```bash
go install github.com/blck-snwmn/dependencylintgo/cmd/llinter@latest
```

## 使い方

### 基本的な使用方法

```bash
llinter ./...
```

### 設定ファイルの指定

```bash
llinter -config=.llinter.yaml ./...
```

## 設定ファイル

`.llinter.yaml`という名前のYAMLファイルをプロジェクトのルートに配置します。

```yaml
rules:
  - path: ["internal/**/*.go"]  # 適用対象のファイルパターン
    deny:
      - "fmt"                   # 禁止するimport
      - "github.com/unauthorized/**"
    allow:
      - "os"                    # 明示的に許可するimport（denyよりも優先される）
  
  - path: ["pkg/**/*.go"]
    deny:
      - "internal/**"           # 内部パッケージを外部に公開しない
```

### ルールの説明

- `path`: ルールを適用するファイルパスのパターン（glob形式）
- `deny`: 禁止するimportパスのパターン
- `allow`: 明示的に許可するimportパスのパターン（denyよりも優先される）

## パターンマッチング

- `*`: 単一ディレクトリ内の任意の文字列にマッチ
- `**`: 複数ディレクトリを横断する任意の文字列にマッチ

## CI統合

### GitHub Actions

```yaml
name: Go Lint

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.24'
    - name: Install LLinter
      run: go install github.com/blck-snwmn/dependencylintgo/cmd/llinter@latest
    - name: Run LLinter
      run: llinter ./...
```

## ライセンス

MIT 