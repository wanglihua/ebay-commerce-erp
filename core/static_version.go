package core

import (
	"github.com/revel/revel"
)

// 这里设定个值，就不从配置文件取了，避免总是去改配置文件
var static_version string = "201804221856"

func GetStaticVersion() string {
	if len(static_version) == 0 {
		static_version = revel.Config.StringDefault("static.version", "1")
	}

	return static_version
}
