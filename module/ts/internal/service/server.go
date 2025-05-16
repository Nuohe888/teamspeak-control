package service

import (
	"easy-fiber-admin/model/ts"
	"easy-fiber-admin/module/ts/internal/tool"
	"easy-fiber-admin/module/ts/internal/utils"
	"easy-fiber-admin/module/ts/internal/vo"
	"easy-fiber-admin/pkg/logger"
	"easy-fiber-admin/pkg/sql"
	"errors"
	"gorm.io/gorm"
)

type ServerSrv struct {
	db  *gorm.DB
	log logger.ILog
}

var serverSrv *ServerSrv

func InitServerSrv() {
	serverSrv = &ServerSrv{
		db:  sql.Get(),
		log: logger.Get(),
	}
}

func GetServerSrv() *ServerSrv {
	if serverSrv == nil {
		panic("service server init failed")
	}
	return serverSrv
}

func (i *ServerSrv) Add(server *ts.Server) error {
	return i.db.Create(&server).Error
}

func (i *ServerSrv) Del(id any) error {
	return i.db.Where("id = ?", id).Delete(&ts.Server{}).Error
}

func (i *ServerSrv) Put(id any, server *ts.Server) error {
	var _server ts.Server
	i.db.Where("id = ?", id).Find(&_server)
	if *_server.Id == 0 {
		return errors.New("不存在该Id")
	}

	utils.MergeStructs(&_server, server)

	return i.db.Save(&_server).Error
}

func (i *ServerSrv) Get(id any) ts.Server {
	var server ts.Server
	i.db.Where("id = ?", id).Find(&server)
	return server
}

func (i *ServerSrv) List(page, pageSize int) *vo.List {
	var items []ts.Server
	var total int64
	if pageSize == 0 {
		pageSize = 20
	}
	db := i.db
	i.db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&items)
	db.Model(&ts.Server{}).Count(&total)
	return &vo.List{
		Items: items,
		Total: total,
	}
}

func (i *ServerSrv) ListAll() []ts.Server {
	var items []ts.Server
	i.db.Find(&items)
	return items
}

func (i *ServerSrv) Check(id any) error {
	var server ts.Server
	i.db.Where("id = ?", id).Find(&server)
	if *server.Id == 0 {
		return errors.New("不存在该Id")
	}

	ts3 := new(tool.Ts)
	ts3.Username = *server.Username
	ts3.Password = *server.Password
	ts3.Host = *server.Host
	ts3.Port = *server.Port

	err := ts3.CheckServer()
	if err != nil {
		return err
	}

	return nil
}
