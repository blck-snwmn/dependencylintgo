package config_test

import (
	"testing"

	"github.com/blck-snwmn/dependencylintgo/analyzer/config"
)

func TestIsFilePathMatched(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		patterns []string
		want     bool
	}{
		{
			name:     "空のパターン",
			filePath: "src/example/main.go",
			patterns: []string{},
			want:     false,
		},
		{
			name:     "完全一致",
			filePath: "src/example/main.go",
			patterns: []string{"src/example/main.go"},
			want:     true,
		},
		{
			name:     "ワイルドカード一致",
			filePath: "src/example/main.go",
			patterns: []string{"src/example/*.go"},
			want:     true,
		},
		{
			name:     "ワイルドカード不一致",
			filePath: "src/example/main.go",
			patterns: []string{"src/other/*.go"},
			want:     false,
		},
		{
			name:     "複数パターン一致",
			filePath: "src/example/main.go",
			patterns: []string{"src/other/*.go", "src/example/*.go"},
			want:     true,
		},
		{
			name:     "再帰的ワイルドカード",
			filePath: "src/example/sub/main.go",
			patterns: []string{"src/example/**/*.go"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.IsFilePathMatched(tt.filePath, tt.patterns)
			if got != tt.want {
				t.Errorf("IsFilePathMatched() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsImportPathMatched(t *testing.T) {
	tests := []struct {
		name       string
		importPath string
		patterns   []string
		want       bool
	}{
		{
			name:       "空のパターン",
			importPath: "fmt",
			patterns:   []string{},
			want:       false,
		},
		{
			name:       "完全一致",
			importPath: "fmt",
			patterns:   []string{"fmt"},
			want:       true,
		},
		{
			name:       "ワイルドカード一致",
			importPath: "github.com/example/pkg",
			patterns:   []string{"github.com/example/*"},
			want:       true,
		},
		{
			name:       "ワイルドカード不一致",
			importPath: "github.com/example/pkg",
			patterns:   []string{"github.com/other/*"},
			want:       false,
		},
		{
			name:       "複数パターン一致",
			importPath: "github.com/example/pkg",
			patterns:   []string{"github.com/other/*", "github.com/example/*"},
			want:       true,
		},
		{
			name:       "サブパッケージ一致",
			importPath: "github.com/example/pkg/sub",
			patterns:   []string{"github.com/example/**"},
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.IsImportPathMatched(tt.importPath, tt.patterns)
			if got != tt.want {
				t.Errorf("IsImportPathMatched() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindMatchingRule(t *testing.T) {
	cfg := &config.Config{
		Rules: []config.Rule{
			{
				Path:  []string{"src/example/*.go"},
				Deny:  []string{"fmt"},
				Allow: []string{"os"},
			},
			{
				Path:  []string{"src/internal/**/*.go"},
				Deny:  []string{"github.com/unauthorized/**"},
				Allow: []string{"github.com/allowed/**"},
			},
		},
	}

	tests := []struct {
		name     string
		filePath string
		wantRule *config.Rule
	}{
		{
			name:     "最初のルールに一致",
			filePath: "src/example/main.go",
			wantRule: &cfg.Rules[0],
		},
		{
			name:     "2番目のルールに一致",
			filePath: "src/internal/pkg/main.go",
			wantRule: &cfg.Rules[1],
		},
		{
			name:     "どのルールにも一致しない",
			filePath: "src/other/main.go",
			wantRule: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := config.FindMatchingRule(cfg, tt.filePath)
			if (got == nil) != (tt.wantRule == nil) {
				t.Errorf("FindMatchingRule() = %v, want %v", got, tt.wantRule)
				return
			}

			if got != nil && tt.wantRule != nil {
				// ポインタの比較ではなく内容の比較
				if len(got.Path) != len(tt.wantRule.Path) ||
					len(got.Deny) != len(tt.wantRule.Deny) ||
					len(got.Allow) != len(tt.wantRule.Allow) {
					t.Errorf("FindMatchingRule() = %v, want %v", got, tt.wantRule)
				}
			}
		})
	}
}
