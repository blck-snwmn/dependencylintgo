package importcheck_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/blck-snwmn/dependencylintgo/analyzer/importcheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

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

	// 設定ファイルのパスをフラグで直接指定
	// 環境変数を使って設定ファイルのパスを渡す
	oldArgs := os.Args
	os.Args = []string{"cmd", "-config=" + configPath}
	defer func() { os.Args = oldArgs }()

	// 設定ファイルのパスをAnalyzerに直接設定
	importcheck.Analyzer.Flags.Set("config", configPath)

	// テスト実行
	analysistest.Run(t, testdata, importcheck.Analyzer, "example")
}
