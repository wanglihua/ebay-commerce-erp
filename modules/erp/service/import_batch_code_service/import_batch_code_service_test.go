package import_batch_code_service

import (
	"testing"
	"fmt"
)

func Test_get_import_batch_code(t *testing.T) {
	var import_code = Get_import_batch_code("12345")
	fmt.Println("import_code", import_code)

	import_code = Get_import_batch_code("abcde")
	fmt.Println("import_code", import_code)
}
