package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/blck-snwmn/dependencylintgo/analyzer/config"
)

func TestLoadConfig(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// テスト用の設定ファイルを作成
	configContent := `
rules:
  - path: ["src/example/*.go"]
    deny:
      - "fmt"
    allow:
      - "os"
  - path: ["src/internal/**/*.go"]
    deny:
      - "github.com/unauthorized/**"
    allow:
      - "github.com/allowed/**"
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

	// 設定内容を検証
	if len(cfg.Rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(cfg.Rules))
	}

	// 1つ目のルールを検証
	if len(cfg.Rules[0].Path) != 1 || cfg.Rules[0].Path[0] != "src/example/*.go" {
		t.Errorf("Unexpected path in first rule: %v", cfg.Rules[0].Path)
	}

	if len(cfg.Rules[0].Deny) != 1 || cfg.Rules[0].Deny[0] != "fmt" {
		t.Errorf("Unexpected deny list in first rule: %v", cfg.Rules[0].Deny)
	}

	if len(cfg.Rules[0].Allow) != 1 || cfg.Rules[0].Allow[0] != "os" {
		t.Errorf("Unexpected allow list in first rule: %v", cfg.Rules[0].Allow)
	}

	// 2つ目のルールを検証
	if len(cfg.Rules[1].Path) != 1 || cfg.Rules[1].Path[0] != "src/internal/**/*.go" {
		t.Errorf("Unexpected path in second rule: %v", cfg.Rules[1].Path)
	}
}

func TestLoadConfigWithWildcards(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// ワイルドカードパターンを含む設定ファイルを作成
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

	// 設定内容を検証
	if len(cfg.Rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(cfg.Rules))
	}

	// 1つ目のルールのdenyリストを検証
	if len(cfg.Rules[0].Deny) != 3 {
		t.Errorf("Expected 3 deny patterns in first rule, got %d", len(cfg.Rules[0].Deny))
	} else {
		// 単純なワイルドカードパターン
		if cfg.Rules[0].Deny[1] != "github.com/unauthorized/*" {
			t.Errorf("Expected simple wildcard pattern, got %s", cfg.Rules[0].Deny[1])
		}

		// 再帰的ワイルドカードパターン
		if cfg.Rules[0].Deny[2] != "github.com/forbidden/**" {
			t.Errorf("Expected recursive wildcard pattern, got %s", cfg.Rules[0].Deny[2])
		}
	}

	// 1つ目のルールのallowリストを検証
	if len(cfg.Rules[0].Allow) != 2 {
		t.Errorf("Expected 2 allow patterns in first rule, got %d", len(cfg.Rules[0].Allow))
	} else {
		// 単純なワイルドカードパターン
		if cfg.Rules[0].Allow[1] != "github.com/allowed/*" {
			t.Errorf("Expected simple wildcard pattern, got %s", cfg.Rules[0].Allow[1])
		}
	}

	// 2つ目のルールのdenyリストを検証
	if len(cfg.Rules[1].Deny) != 1 {
		t.Errorf("Expected 1 deny pattern in second rule, got %d", len(cfg.Rules[1].Deny))
	} else {
		// 再帰的ワイルドカードパターン
		if cfg.Rules[1].Deny[0] != "internal/**" {
			t.Errorf("Expected recursive wildcard pattern, got %s", cfg.Rules[1].Deny[0])
		}
	}

	// 2つ目のルールのallowリストを検証
	if len(cfg.Rules[1].Allow) != 3 {
		t.Errorf("Expected 3 allow patterns in second rule, got %d", len(cfg.Rules[1].Allow))
	} else {
		// 単純なワイルドカードパターン
		if cfg.Rules[1].Allow[1] != "golang.org/x/*" {
			t.Errorf("Expected simple wildcard pattern, got %s", cfg.Rules[1].Allow[1])
		}

		// 再帰的ワイルドカードパターン
		if cfg.Rules[1].Allow[2] != "github.com/approved/**" {
			t.Errorf("Expected recursive wildcard pattern, got %s", cfg.Rules[1].Allow[2])
		}
	}

	// ロードした設定を使ってパターンマッチングをテスト
	// 単純なワイルドカードパターンのテスト
	rule := &cfg.Rules[0]
	if !config.IsImportPathMatched("github.com/unauthorized/pkg", rule.Deny) {
		t.Error("Simple wildcard pattern should match github.com/unauthorized/pkg")
	}

	if config.IsImportPathMatched("github.com/unauthorized/sub/pkg", rule.Deny) {
		t.Error("Simple wildcard pattern should not match github.com/unauthorized/sub/pkg")
	}

	// 再帰的ワイルドカードパターンのテスト
	if !config.IsImportPathMatched("github.com/forbidden/pkg", rule.Deny) {
		t.Error("Recursive wildcard pattern should match github.com/forbidden/pkg")
	}

	if !config.IsImportPathMatched("github.com/forbidden/sub/pkg", rule.Deny) {
		t.Error("Recursive wildcard pattern should match github.com/forbidden/sub/pkg")
	}
}
