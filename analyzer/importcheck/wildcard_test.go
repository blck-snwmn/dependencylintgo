package importcheck_test

import (
	"testing"

	"github.com/blck-snwmn/dependencylintgo/analyzer/config"
)

func TestIsImportPathMatchedWithWildcard(t *testing.T) {
	tests := []struct {
		name       string
		importPath string
		patterns   []string
		want       bool
	}{
		{
			name:       "単純なワイルドカード一致",
			importPath: "github.com/unauthorized/pkg",
			patterns:   []string{"github.com/unauthorized/*"},
			want:       true,
		},
		{
			name:       "単純なワイルドカード不一致",
			importPath: "github.com/unauthorized/sub/pkg",
			patterns:   []string{"github.com/unauthorized/*"},
			want:       false,
		},
		{
			name:       "再帰的ワイルドカード一致（直下）",
			importPath: "github.com/forbidden/pkg",
			patterns:   []string{"github.com/forbidden/**"},
			want:       true,
		},
		{
			name:       "再帰的ワイルドカード一致（サブディレクトリ）",
			importPath: "github.com/forbidden/sub/pkg",
			patterns:   []string{"github.com/forbidden/**"},
			want:       true,
		},
		{
			name:       "再帰的ワイルドカード一致（深いサブディレクトリ）",
			importPath: "github.com/forbidden/sub/deep/pkg",
			patterns:   []string{"github.com/forbidden/**"},
			want:       true,
		},
		{
			name:       "複数パターン（一致するものがある）",
			importPath: "github.com/allowed/pkg",
			patterns:   []string{"github.com/unauthorized/*", "github.com/allowed/*"},
			want:       true,
		},
		{
			name:       "複数パターン（一致するものがない）",
			importPath: "github.com/other/pkg",
			patterns:   []string{"github.com/unauthorized/*", "github.com/allowed/*"},
			want:       false,
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
