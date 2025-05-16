package ts

import (
	"easy-fiber-admin/module/ts/internal/controller"
	"easy-fiber-admin/module/ts/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(r fiber.Router) {
	auth := r.Group("ts")
	auth.Use(middleware.JWT()).
		Use(middleware.Casbin())

	auth.Post("server", controller.ServerCtl.Add)
	auth.Delete("server/:id", controller.ServerCtl.Del)
	auth.Put("server/:id", controller.ServerCtl.Put)
	auth.Get("server", controller.ServerCtl.Get)
	auth.Get("server/list", controller.ServerCtl.List)
	auth.Get("server/list/all", controller.ServerCtl.ListAll)
	auth.Get("server/check/:id", controller.ServerCtl.Check)

	auth.Post("order", controller.OrderCtl.Add)
	auth.Delete("order/:uuid", controller.OrderCtl.Del)
	auth.Put("order/:uuid", controller.OrderCtl.Put)
	auth.Get("order", controller.OrderCtl.Get)
	auth.Get("order/list", controller.OrderCtl.List)
	auth.Get("order/list/all", controller.OrderCtl.ListAll)
	auth.Post("order/tsrun", controller.OrderCtl.Runts)
	auth.Post("order/tsdel", controller.OrderCtl.DelTs)
	auth.Post("order/tsstatus", controller.OrderCtl.TsStatus)
	auth.Post("order/tsinfo", controller.OrderCtl.TsInfo)
}
