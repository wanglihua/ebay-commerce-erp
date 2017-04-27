package import_batch_code_service

import (
	"time"
)

func Get_import_batch_code() int64 {
	var time_now = time.Now()

	var import_batch_code int64 = int64(time_now.Year())*10000000000 + int64(time_now.Month())*100000000 + int64(time_now.Day()*1000000) + int64(time_now.Hour()*10000) + int64(time_now.Minute()*100) + int64(time_now.Second())

	return import_batch_code
}
