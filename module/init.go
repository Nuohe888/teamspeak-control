package module

import (
	"easy-fiber-admin/module/system"
	"easy-fiber-admin/module/ts"
)

func Init() {
	system.Init()
	ts.Init()
}
