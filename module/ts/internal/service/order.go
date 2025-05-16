package service

import (
	"easy-fiber-admin/model/ts"
	"easy-fiber-admin/module/ts/internal/tool"
	"easy-fiber-admin/module/ts/internal/utils"
	"easy-fiber-admin/module/ts/internal/vo"
	"easy-fiber-admin/pkg/logger"
	"easy-fiber-admin/pkg/sql"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderSrv struct {
	db  *gorm.DB
	log logger.ILog
}

var orderSrv *OrderSrv

func InitOrderSrv() {
	orderSrv = &OrderSrv{
		db:  sql.Get(),
		log: logger.Get(),
	}
}

func GetOrderSrv() *OrderSrv {
	if orderSrv == nil {
		panic("service order init failed")
	}
	return orderSrv
}

func (i *OrderSrv) Add(order *ts.Order) error {
	if order.Uuid != nil {
		return errors.New("禁止传入UUID")
	}
	_uuid := uuid.NewString()
	order.Uuid = &_uuid
	return i.db.Create(&order).Error
}

func (i *OrderSrv) Del(uuid any) error {
	return i.db.Where("uuid = ?", uuid).Delete(&ts.Order{}).Error
}

func (i *OrderSrv) Put(uuid any, order *ts.Order) error {
	var _order ts.Order
	i.db.Where("uuid = ?", uuid).Find(&_order)
	if *_order.Uuid == "" {
		return errors.New("不存在该Uuid")
	}

	utils.MergeStructs(&_order, order)

	return i.db.Save(&_order).Error
}

func (i *OrderSrv) Get(uuid any) ts.Order {
	var order ts.Order
	i.db.Where("uuid = ?", uuid).Find(&order)
	return order
}

func (i *OrderSrv) List(page, pageSize int) *vo.List {
	var items []ts.Order
	var total int64
	if pageSize == 0 {
		pageSize = 20
	}
	db := i.db
	i.db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&items)
	db.Model(&ts.Order{}).Count(&total)
	return &vo.List{
		Items: items,
		Total: total,
	}
}

func (i *OrderSrv) ListAll() []ts.Order {
	var items []ts.Order
	i.db.Find(&items)
	return items
}

func (i *OrderSrv) RunTs(uuid any) error {
	var order ts.Order
	i.db.Where("uuid = ?", uuid).Find(&order)
	if *order.Uuid == "" {
		return errors.New("不存在该Uuid")
	}
	var server ts.Server
	i.db.Where("id = ?", order.ServerId).Find(&server)
	if *server.Id == 0 {
		return errors.New("不存在该ServerId")
	}
	ts3 := new(tool.Ts)
	ts3.Username = *server.Username
	ts3.Password = *server.Password
	ts3.Host = *server.Host
	ts3.Port = *server.Port
	ts3.Uuid = *order.Uuid
	ts3.ImageName = *server.ImageName
	ts3.DefaultVoicePort = *order.TsDefaultVoicePort
	ts3.QueryPort = *order.TsQueryPort
	ts3.FiletransferPort = *order.TsFiletransferPort

	err := ts3.Run()
	if err != nil {
		return err
	}

	logMap, err := ts3.ParseLog()
	if err != nil {
		return err
	}

	var (
		token     = logMap["token"]
		loginname = logMap["loginname"]
		password  = logMap["password"]
		apikey    = logMap["apikey"]
	)

	order.TsToken = &token
	order.TsLoginName = &loginname
	order.TsPassword = &password
	order.TsApikey = &apikey

	return i.db.Save(&order).Error
}

func (i *OrderSrv) DelTs(uuid any) error {
	var order ts.Order
	i.db.Where("uuid = ?", uuid).Find(&order)
	if *order.Uuid == "" {
		return errors.New("不存在该Uuid")
	}
	var server ts.Server
	i.db.Where("id = ?", order.ServerId).Find(&server)
	if *server.Id == 0 {
		return errors.New("不存在该ServerId")
	}
	ts3 := new(tool.Ts)
	ts3.Username = *server.Username
	ts3.Password = *server.Password
	ts3.Host = *server.Host
	ts3.Port = *server.Port
	ts3.Uuid = *order.Uuid
	return ts3.Del()
}

func (i *OrderSrv) TsStatus(uuid any) (*vo.TsStatusRes, error) {
	var order ts.Order
	i.db.Where("uuid = ?", uuid).Find(&order)
	if *order.Uuid == "" {
		return nil, errors.New("不存在该Uuid")
	}
	var server ts.Server
	i.db.Where("id = ?", order.ServerId).Find(&server)
	if *server.Id == 0 {
		return nil, errors.New("不存在该ServerId")
	}
	ts3 := new(tool.Ts)
	ts3.Username = *server.Username
	ts3.Password = *server.Password
	ts3.Host = *server.Host
	ts3.Port = *server.Port
	ts3.Uuid = *order.Uuid
	return &vo.TsStatusRes{Status: ts3.Status()}, nil
}

func (i *OrderSrv) TsInfo(uuid any) (*vo.TsInfoRes, error) {
	var order ts.Order
	i.db.Where("uuid = ?", uuid).Find(&order)
	if *order.Uuid == "" {
		return nil, errors.New("不存在该Uuid")
	}
	var server ts.Server
	i.db.Where("id = ?", order.ServerId).Find(&server)
	if *server.Id == 0 {
		return nil, errors.New("不存在该ServerId")
	}
	return &vo.TsInfoRes{
		Domain:           *server.Domain,
		DefaultVoicePort: *order.TsDefaultVoicePort,
		QueryPort:        *order.TsQueryPort,
		FiletransferPort: *order.TsFiletransferPort,
	}, nil
}
