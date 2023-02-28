package data

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
	SendAddress
	ReceiverAddress
	gorm.Model
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

// PlatformTrade 平台订单
type PlatformTrade struct {
	Tid  string `json:"tid"`
	Json string `json:"json"`
	gorm.Model
}

// SendAddress 发货地址
type SendAddress struct {
	SendProvince string `json:"sendProvince"`
	SendCity     string `json:"sendCity"`
	SendDistinct string `json:"sendDistinct"`
	SendTown     string `json:"sendTown"`
	SendDetail   string `json:"sendDetail"`
	SendCode     string `json:"sendCode"`
}

// ReceiverAddress 收货地址
type ReceiverAddress struct {
	ReceiverName     string `json:"receiverName"`
	ReceiverProvince string `json:"receiverProvince"`
	ReceiverCity     string `json:"receiverCity"`
	ReceiverDistinct string `json:"receiverDistinct"`
	ReceiverTown     string `json:"receiverTown"`
	ReceiverDetail   string `json:"receiverDetail"`
	ReceiverCode     string `json:"receiverCode"`
}

// Unit 一对多的关系单元
type Unit struct {
	UnitId   uint   `json:"unitId"`   //一对多的单元id
	UnitName string `json:"unitName"` //一对多的单元名称
}
