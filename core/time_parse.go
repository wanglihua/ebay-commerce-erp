package core

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Parse_time(str string) (time.Time, error) {
	var failed_time = time.Time{}
	var failed_error = errors.New(fmt.Sprintf("%s 不是有效的时间格式！", str))

	//if strings.TrimSpace(str) == "" {
	//	return time.Time{}, nil
	//}

	var date_time_buffer bytes.Buffer

	var str_list = strings.Split(str, " ")
	if len(str_list) >= 1 {
		var date_str_list = strings.Split(str_list[0], "-")
		if len(date_str_list) != 3 {
			return failed_time, failed_error
		}

		year, err := strconv.ParseInt(strings.TrimPrefix(date_str_list[0], "0"), 10, 64)
		if err != nil {
			return failed_time, failed_error
		}

		month, err := strconv.ParseInt(strings.TrimPrefix(date_str_list[1], "0"), 10, 64)
		if err != nil {
			return failed_time, failed_error
		}

		day, err := strconv.ParseInt(strings.TrimPrefix(date_str_list[2], "0"), 10, 64)
		if err != nil {
			return failed_time, failed_error
		}

		date_time_buffer.WriteString(fmt.Sprintf("%04d", year))
		date_time_buffer.WriteString("-")
		date_time_buffer.WriteString(fmt.Sprintf("%02d", month))
		date_time_buffer.WriteString("-")
		date_time_buffer.WriteString(fmt.Sprintf("%02d", day))
	}

	if len(str_list) == 1 {
		date_time_buffer.WriteString(" 00:00:00") // 这里面的开头的空格不能少
	}

	if len(str_list) == 2 {
		date_time_buffer.WriteString(" ")

		var time_str_list = strings.Split(str_list[1], ":")
		if len(time_str_list) > 3 {
			return failed_time, failed_error
		}

		var hour int64 = 0
		var err error
		if len(time_str_list) >= 1 {
			hour, err = strconv.ParseInt(strings.TrimPrefix(time_str_list[0], "0"), 10, 64)
			if err != nil {
				return failed_time, failed_error
			}
		}

		var minute int64 = 0
		if len(time_str_list) >= 2 {
			minute, err = strconv.ParseInt(strings.TrimPrefix(time_str_list[1], "0"), 10, 64)
			if err != nil {
				return failed_time, failed_error
			}
		}

		var second int64 = 0
		if len(time_str_list) == 3 {
			second, err = strconv.ParseInt(strings.TrimPrefix(time_str_list[2], "0"), 10, 64)
			if err != nil {
				return failed_time, failed_error
			}
		}

		date_time_buffer.WriteString(fmt.Sprintf("%02d", hour))
		date_time_buffer.WriteString(":")
		date_time_buffer.WriteString(fmt.Sprintf("%02d", minute))
		date_time_buffer.WriteString(":")
		date_time_buffer.WriteString(fmt.Sprintf("%02d", second))
	}

	if len(str_list) > 2 {
		return failed_time, failed_error
	}

	var date_time_str = date_time_buffer.String()

	if date_time_str == "0001-01-01 00:00:00" {
		return time.Time{}, nil
	}

	r, err := time.ParseInLocation("2006-01-02 15:04:05", date_time_str, time.Local)

	if err == nil {
		return r, nil
	}

	return failed_time, failed_error
}
