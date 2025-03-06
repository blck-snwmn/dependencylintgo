package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config は設定ファイルの構造体だ
type Config struct {
	Rules []Rule `yaml:"rules"`
}

// Rule はimportルールを定義するだ
type Rule struct {
	Path  []string `yaml:"path"`  // 適用するファイルパスパターン
	Deny  []string `yaml:"deny"`  // 禁止するimportパターン
	Allow []string `yaml:"allow"` // 許可するimportパターン（denyよりも優先される）
}

// LoadConfig は設定ファイルを読み込むだ
func LoadConfig(configPath string) (*Config, error) {
	// 指定されたパスが相対パスなら絶対パスに変換
	if !filepath.IsAbs(configPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		configPath = filepath.Join(cwd, configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
