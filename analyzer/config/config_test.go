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
