package order_logistics_service

import (
	"bytes"
	"fmt"
)

type match_trace struct {
	match_debug bytes.Buffer
}

func (mt match_trace) Println(a ...interface{}) {
	mt.match_debug.WriteString("<p>")
	mt.match_debug.WriteString(fmt.Sprint(a))
	mt.match_debug.WriteString("</p>")
	mt.match_debug.WriteString(fmt.Sprintln())
}

func (mt match_trace) String(a ...interface{}) string {
	return mt.match_debug.String()
}
