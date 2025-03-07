package example

import (
	"os"
)

// InternalTest は内部パッケージのimportをテストする関数
func InternalTest() {
	os.Exit(0)
	// 実際には存在しないパッケージなので、コンパイルエラーを避けるためにコメントアウト
	// _ = foo.Dummy
}
