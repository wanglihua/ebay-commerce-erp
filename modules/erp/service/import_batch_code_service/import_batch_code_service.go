package import_batch_code_service

import (
	"time"
	"bytes"
	"fmt"
	"strconv"
)

func Get_import_batch_code(login_name string) int64 {
	var time_now = time.Now()

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%02d", time_now.Year()%100))
	buffer.WriteString(fmt.Sprintf("%02d", time_now.Month()))
	buffer.WriteString(fmt.Sprintf("%02d", time_now.Day()%100))
	buffer.WriteString(fmt.Sprintf("%02d", time_now.Hour()))
	buffer.WriteString(fmt.Sprintf("%02d", time_now.Minute()))

	work_num, err := strconv.ParseInt(login_name, 10, 64)
	if err != nil {
		work_num = 88888
	}
	buffer.WriteString(fmt.Sprintf("%05d", work_num))

	//var import_batch_code int64 = int64(time_now.Year())*10000000000 + int64(time_now.Month())*100000000 + int64(time_now.Day()*1000000) + int64(time_now.Hour()*10000) + int64(time_now.Minute()*100) + int64(time_now.Second())
	import_batch_code, _ := strconv.ParseInt(buffer.String(), 10, 64)

	return import_batch_code
}
