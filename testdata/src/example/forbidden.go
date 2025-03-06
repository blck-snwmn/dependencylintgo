package example

import (
	"fmt" // want "import \"fmt\" is not allowed in this file based on configuration"
	"os"
)

func Example() {
	fmt.Println("This is a test")
	os.Exit(0)
}
