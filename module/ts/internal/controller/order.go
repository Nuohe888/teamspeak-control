package controller

import (
	"easy-fiber-admin/model/ts"
	"easy-fiber-admin/module/ts/internal/service"
	"easy-fiber-admin/module/ts/internal/vo"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type orderCtl struct {
	srv *service.OrderSrv
}

var OrderCtl *orderCtl

func InitOrderCtl() {
	OrderCtl = &orderCtl{
		srv: service.GetOrderSrv(),
	}
}

func (i *orderCtl) Add(c *fiber.Ctx) error {
	var req ts.Order
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	if err := i.srv.Add(&req); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *orderCtl) Del(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	if err := i.srv.Del(uuid); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *orderCtl) Put(c *fiber.Ctx) error {
	uuid := c.Params("uuid")
	var req ts.Order
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	if err := i.srv.Put(uuid, &req); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *orderCtl) Get(c *fiber.Ctx) error {
	uuid := c.Query("uuid")
	return vo.ResultOK(i.srv.Get(uuid), c)
}

func (i *orderCtl) List(c *fiber.Ctx) error {
	page := c.Query("page")
	pageSize := c.Query("pageSize")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	return vo.ResultOK(i.srv.List(pageInt, pageSizeInt), c)
}

func (i *orderCtl) ListAll(c *fiber.Ctx) error {
	return vo.ResultOK(i.srv.ListAll(), c)
}

func (i *orderCtl) Runts(c *fiber.Ctx) error {
	var req vo.TsUuidReq
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	if err := i.srv.RunTs(req.Uuid); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *orderCtl) DelTs(c *fiber.Ctx) error {
	var req vo.TsUuidReq
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	if err := i.srv.DelTs(req.Uuid); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *orderCtl) TsStatus(c *fiber.Ctx) error {
	var req vo.TsUuidReq
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	res, err := i.srv.TsStatus(req.Uuid)
	if err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(res, c)
}

func (i *orderCtl) TsInfo(c *fiber.Ctx) error {
	var req vo.TsUuidReq
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	res, err := i.srv.TsInfo(req.Uuid)
	if err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(res, c)
}
