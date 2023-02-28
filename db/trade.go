package db

import "gorm.io/gorm"

// TradeMain 订单主要字段
type TradeMain struct {
	Sid           uint   `json:"sid"`           // 订单号分布式id
	Tid           string `json:"tid"`           // 平台id
	CompanyId     uint   `json:"CompanyId"`     // 公司id
	WarehouseId   uint   `json:"warehouseId"`   // 仓库id
	WaybillNumber string `json:"waybillNumber"` // 运单号
	TemplateId    uint   `json:"templateId"`    // 模板id
	SysStatus     string `json:"sysStatus"`     // 系统状态
	CreateMethod  string `json:"createMethod"`  // 创建方式 sys系统手工建 platform_sync同步平台
	PlatformCode  string `json:"platformCode"`  // 平台编码
	PlatformName  string `json:"platformName"`  // 平台名称
	Remark        string `json:"remark"`        // 备注
	gorm.Model
}

// Unit 一对多的关系单元
type Unit struct {
	UnitId   uint   `json:"unitId"`   //一对多的单元id
	UnitName string `json:"unitName"` //一对多的单元名称
}

// TradeHasMany 订单一对多关系表
type TradeHasMany struct {
	Sid       uint   `json:"sid"`       // 订单号分布式id
	CompanyId uint   `json:"CompanyId"` // 公司id
	Type      int    `json:"type"`      // 类型
	TypeName  string `json:"typeName"`  // 类型名称
	Unit
	gorm.Model
}
