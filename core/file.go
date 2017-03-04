package core

import (
	"os"
	"github.com/revel/revel"
)

func Force_write_to_file(file_name string, content string) {

	is_file_exist := func(file_name string) bool {
		var exist = true
		if _, err := os.Stat(file_name); os.IsNotExist(err) {
			exist = false
		}
		return exist
	}

	if is_file_exist(file_name) {
		err := os.Remove(file_name)
		HandleError(err)
	}

	var file *os.File
	var err error

	file, err = os.Create(file_name) //创建文件
	if err != nil {
		revel.ERROR.Println(err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		revel.ERROR.Println(err)
	}

	file.Sync()
	file.Close()
}
