package ts

import (
	"easy-fiber-admin/module/ts/internal/controller"
	"easy-fiber-admin/module/ts/internal/service"
)

func Init() {
	service.Init()
	controller.Init()
}
