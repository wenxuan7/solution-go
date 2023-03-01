package data

import "gorm.io/gorm"

// TradeMain 订单主要字段
type TradeMain struct {
	Sid           uint   // 订单号分布式id
	Tid           string // 平台id
	CompanyId     uint   // 公司id
	WarehouseId   uint   // 仓库id
	WaybillNumber string // 运单号
	TemplateId    uint   // 模板id
	SysStatus     string // 系统状态
	CreateMethod  string // 创建方式 sys系统手工建 platform_sync同步平台
	PlatformCode  string // 平台编码
	PlatformName  string // 平台名称
	Remark        string // 备注
	SendAddress
	ReceiverAddress
	gorm.Model
}

func (tm *TradeMain) TableName() string {
	return "trade_main"
}

// TradeHasMany 订单一对多关系表
type TradeHasMany struct {
	Sid       uint   // 订单号分布式id
	CompanyId uint   // 公司id
	Type      int    // 类型
	TypeName  string // 类型名称
	Unit
	gorm.Model
}

func (th *TradeHasMany) TableName() string {
	return "trade_has_many"
}

// PlatformTrade 平台订单
type PlatformTrade struct {
	Tid  string
	Json string
	gorm.Model
}

func (pt *PlatformTrade) TableName() string {
	return "platform_trade"
}

// SendAddress 发货地址
type SendAddress struct {
	SendProvince string
	SendCity     string
	SendDistinct string
	SendTown     string
	SendDetail   string
	SendCode     string
}

// ReceiverAddress 收货地址
type ReceiverAddress struct {
	ReceiverName     string
	ReceiverProvince string
	ReceiverCity     string
	ReceiverDistinct string
	ReceiverTown     string
	ReceiverDetail   string
	ReceiverCode     string
}

// Unit 一对多的关系单元
type Unit struct {
	UnitId   uint   //一对多的单元id
	UnitName string //一对多的单元名称
}
