package example

import (
	"fmt" // want "import \"fmt\" is not allowed in this file based on configuration"
	"os"
)

// WildcardTest はワイルドカードパターンのテスト用関数
func WildcardTest() {
	fmt.Println("This is a wildcard test")
	os.Exit(0)
}
