package order_logistics_service

import (
	"bytes"
	"fmt"
)

type match_trace struct {
	is_trace    bool
	match_debug bytes.Buffer
}

func (mt *match_trace) set_is_trace(setting int64) {
	if setting > 0 {
		mt.is_trace = true
	} else {
		mt.is_trace = false
	}
}

func (mt *match_trace) Println(a ...interface{}) {
	if mt.is_trace == true {
		//mt.match_debug.WriteString("<p>")
		mt.match_debug.WriteString(fmt.Sprint(a...))
		//mt.match_debug.WriteString("</p>")
		mt.match_debug.WriteString(fmt.Sprintln())
	}
}

func (mt *match_trace) PrintlnIndent(indent_level int64, a ...interface{}) {
	if mt.is_trace == true {
		//mt.match_debug.WriteString(fmt.Sprintf("<p style='text-indent:%dem'>", indent_level*4)) // 一次缩进 4 个字符
		for i := 0; i < int(indent_level*4); i++ {
			mt.match_debug.WriteString(fmt.Sprint(" ")) // 一次缩进 4 个字符
		}
		mt.match_debug.WriteString(fmt.Sprint(a...))
		//mt.match_debug.WriteString("</p>")
		mt.match_debug.WriteString(fmt.Sprintln())
	}
}

func (mt *match_trace) String() string {
	return mt.match_debug.String()
}
