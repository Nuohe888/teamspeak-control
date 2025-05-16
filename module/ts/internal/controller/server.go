package controller

import (
	"easy-fiber-admin/model/ts"
	"easy-fiber-admin/module/ts/internal/service"
	"easy-fiber-admin/module/ts/internal/vo"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type serverCtl struct {
	srv *service.ServerSrv
}

var ServerCtl *serverCtl

func InitServerCtl() {
	ServerCtl = &serverCtl{
		srv: service.GetServerSrv(),
	}
}

func (i *serverCtl) Add(c *fiber.Ctx) error {
	var req ts.Server
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	if err := i.srv.Add(&req); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *serverCtl) Del(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := i.srv.Del(id); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *serverCtl) Put(c *fiber.Ctx) error {
	id := c.Params("id")
	var req ts.Server
	if err := vo.BodyParser(&req, c); err != nil {
		return err
	}
	if err := i.srv.Put(id, &req); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}

func (i *serverCtl) Get(c *fiber.Ctx) error {
	id := c.Query("id")
	return vo.ResultOK(i.srv.Get(id), c)
}

func (i *serverCtl) List(c *fiber.Ctx) error {
	page := c.Query("page")
	pageSize := c.Query("pageSize")
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	return vo.ResultOK(i.srv.List(pageInt, pageSizeInt), c)
}

func (i *serverCtl) ListAll(c *fiber.Ctx) error {
	return vo.ResultOK(i.srv.ListAll(), c)
}

func (i *serverCtl) Check(c *fiber.Ctx) error {
	id := c.Query("id")
	if err := i.srv.Check(id); err != nil {
		return vo.ResultErr(err, c)
	}
	return vo.ResultOK(nil, c)
}
