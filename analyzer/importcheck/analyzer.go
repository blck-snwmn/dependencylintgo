package importcheck

import (
	"go/ast"
	"strings"

	"github.com/blck-snwmn/dependencylintgo/analyzer/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var configFile string

// Analyzer はimportチェック用のanalyzerだ
var Analyzer = &analysis.Analyzer{
	Name: "importcheck",
	Doc:  "checks for disallowed imports based on configuration",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func init() {
	Analyzer.Flags.StringVar(&configFile, "config", ".llinter.yaml", "configuration file path")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// 設定ファイルの読み込み
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		// 設定ファイルが見つからない場合はスキップする
		if strings.Contains(err.Error(), "no such file") {
			return nil, nil
		}
		return nil, err
	}

	nodeFilter := []ast.Node{
		(*ast.ImportSpec)(nil),
	}

	inspector.Preorder(nodeFilter, func(n ast.Node) {
		importSpec := n.(*ast.ImportSpec)
		importPath := strings.Trim(importSpec.Path.Value, "\"")

		// ファイルパスを取得
		filePath := pass.Fset.Position(n.Pos()).Filename

		// ファイルパスを処理
		relPath := filePath

		// テスト環境では、ファイルパスにtestdata/src/が含まれる
		if strings.Contains(filePath, "testdata/src/") {
			// testdata/src/ 以降のパスを抽出
			parts := strings.Split(filePath, "testdata/src/")
			if len(parts) > 1 {
				relPath = parts[1]
			}
		} else if strings.Contains(filePath, "/src/") {
			// 通常の Go パッケージの場合は /src/ 以降を抽出
			parts := strings.Split(filePath, "/src/")
			if len(parts) > 1 {
				relPath = "src/" + parts[1]
			}
		}

		// ルールを検索
		rule := config.FindMatchingRule(cfg, relPath)
		if rule == nil {
			return // マッチするルールがなければチェックしない
		}

		// denyリストに含まれているかチェック
		if config.IsImportPathMatched(importPath, rule.Deny) {
			// allowリストで明示的に許可されているか確認
			if !config.IsImportPathMatched(importPath, rule.Allow) {
				// エラー報告
				pass.Reportf(importSpec.Pos(), "import %q is not allowed in this file based on configuration", importPath)
			}
		}
	})

	return nil, nil
}
