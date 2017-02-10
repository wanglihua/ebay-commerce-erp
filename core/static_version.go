package core

import (
	"github.com/revel/revel"
)

var static_version string = ""

func GetStaticVersion() string {
	if len(static_version) == 0 {
		static_version = revel.Config.StringDefault("static.version", "1")
	}

	return static_version
}
