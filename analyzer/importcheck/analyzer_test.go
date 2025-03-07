package importcheck_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/blck-snwmn/dependencylintgo/analyzer/config"
	"github.com/blck-snwmn/dependencylintgo/analyzer/importcheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer は実際のGoファイルに対してAnalyzerを実行するテスト
func TestAnalyzer(t *testing.T) {
	// テストデータのディレクトリを取得
	// プロジェクトルートからの相対パスで指定
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// analyzer/importcheckからプロジェクトルートへ移動
	rootDir := filepath.Join(wd, "../..")
	testdata := filepath.Join(rootDir, "testdata")

	// テスト実行前に設定ファイルのパスを設定
	configPath := filepath.Join(testdata, ".llinter.yaml")

	// 設定ファイルのパスをAnalyzerに直接設定
	importcheck.Analyzer.Flags.Set("config", configPath)

	// テスト実行
	analysistest.Run(t, testdata, importcheck.Analyzer, "example")
}

// TestWildcardPatterns はワイルドカードパターンのマッチング機能をテスト
func TestWildcardPatterns(t *testing.T) {
	// ワイルドカードパターンを含む設定ファイルを作成
	tempDir := t.TempDir()
	configContent := `
rules:
  - path: ["example/forbidden.go"]
    deny:
      - "fmt"
      - "github.com/unauthorized/*"  # 単純なワイルドカード
      - "github.com/forbidden/**"    # 再帰的ワイルドカード
    allow:
      - "os"
      - "github.com/allowed/*"       # 単純なワイルドカード
  - path: ["example/*.go"]
    deny:
      - "internal/**"                # 再帰的ワイルドカード
    allow:
      - "os"
      - "golang.org/x/*"             # 単純なワイルドカード
      - "github.com/approved/**"     # 再帰的ワイルドカード
`
	configPath := filepath.Join(tempDir, ".llinter.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// 設定ファイルを読み込み
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// テストケース
	testCases := []struct {
		name       string
		importPath string
		filePath   string
		wantDenied bool
	}{
		{
			name:       "fmt in forbidden.go",
			importPath: "fmt",
			filePath:   "example/forbidden.go",
			wantDenied: true,
		},
		{
			name:       "os in forbidden.go",
			importPath: "os",
			filePath:   "example/forbidden.go",
			wantDenied: false, // allowリストに含まれている
		},
		{
			name:       "github.com/unauthorized/pkg in forbidden.go",
			importPath: "github.com/unauthorized/pkg",
			filePath:   "example/forbidden.go",
			wantDenied: true, // 単純なワイルドカードに一致
		},
		{
			name:       "github.com/unauthorized/sub/pkg in forbidden.go",
			importPath: "github.com/unauthorized/sub/pkg",
			filePath:   "example/forbidden.go",
			wantDenied: false, // 単純なワイルドカードに一致しない
		},
		{
			name:       "github.com/forbidden/pkg in forbidden.go",
			importPath: "github.com/forbidden/pkg",
			filePath:   "example/forbidden.go",
			wantDenied: true, // 再帰的ワイルドカードに一致
		},
		{
			name:       "github.com/forbidden/sub/pkg in forbidden.go",
			importPath: "github.com/forbidden/sub/pkg",
			filePath:   "example/forbidden.go",
			wantDenied: true, // 再帰的ワイルドカードに一致
		},
		{
			name:       "github.com/allowed/pkg in forbidden.go",
			importPath: "github.com/allowed/pkg",
			filePath:   "example/forbidden.go",
			wantDenied: false, // allowリストに含まれている
		},
		{
			name:       "internal/foo in other.go",
			importPath: "internal/foo",
			filePath:   "example/other.go",
			wantDenied: true, // 再帰的ワイルドカードに一致
		},
		{
			name:       "internal/sub/pkg in other.go",
			importPath: "internal/sub/pkg",
			filePath:   "example/other.go",
			wantDenied: true, // 再帰的ワイルドカードに一致
		},
		{
			name:       "golang.org/x/tools in other.go",
			importPath: "golang.org/x/tools",
			filePath:   "example/other.go",
			wantDenied: false, // allowリストに含まれている
		},
		{
			name:       "github.com/approved/pkg in other.go",
			importPath: "github.com/approved/pkg",
			filePath:   "example/other.go",
			wantDenied: false, // allowリストに含まれている
		},
		{
			name:       "github.com/approved/sub/pkg in other.go",
			importPath: "github.com/approved/sub/pkg",
			filePath:   "example/other.go",
			wantDenied: false, // allowリストに含まれている
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ルールを検索
			rule := config.FindMatchingRule(cfg, tc.filePath)
			if rule == nil {
				t.Fatalf("No matching rule found for %s", tc.filePath)
			}

			// denyリストに含まれているかチェック
			isDenied := config.IsImportPathMatched(tc.importPath, rule.Deny)

			// allowリストで明示的に許可されているか確認
			if isDenied {
				isAllowed := config.IsImportPathMatched(tc.importPath, rule.Allow)
				isDenied = !isAllowed
			}

			if isDenied != tc.wantDenied {
				t.Errorf("Import %s in %s: got denied=%v, want denied=%v", tc.importPath, tc.filePath, isDenied, tc.wantDenied)
			}
		})
	}
}
