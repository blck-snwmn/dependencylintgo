package config

import (
	"path/filepath"
	"strings"
)

// IsFilePathMatched はファイルパスがパターンにマッチするか確認する
func IsFilePathMatched(filePath string, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}

	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, filePath)
		if err == nil && matched {
			return true
		}

		// より複雑なグロブパターン（**など）のサポート
		if strings.Contains(pattern, "**") {
			parts := strings.Split(pattern, "**")
			if len(parts) == 2 {
				if strings.HasPrefix(filePath, parts[0]) && strings.HasSuffix(filePath, parts[1]) {
					return true
				}
			}
		}
	}

	return false
}

// IsImportPathMatched はインポートパスがパターンにマッチするか確認する
func IsImportPathMatched(importPath string, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}

	for _, pattern := range patterns {
		// 完全一致のケース
		if pattern == importPath {
			return true
		}

		// ワイルドカードを含むパターン
		if strings.Contains(pattern, "*") {
			matched, err := filepath.Match(pattern, importPath)
			if err == nil && matched {
				return true
			}
		}

		// プレフィックスマッチ（サブパッケージ含む）
		if strings.HasSuffix(pattern, "/**") {
			prefix := strings.TrimSuffix(pattern, "/**")
			if strings.HasPrefix(importPath, prefix) {
				return true
			}
		}
	}

	return false
}

// FindMatchingRule はファイルパスに適用するルールを見つける
func FindMatchingRule(config *Config, filePath string) *Rule {
	if config == nil || len(config.Rules) == 0 {
		return nil
	}

	for _, rule := range config.Rules {
		if IsFilePathMatched(filePath, rule.Path) {
			return &rule
		}
	}

	return nil
}
