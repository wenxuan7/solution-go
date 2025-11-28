package order

import (
	"github.com/wenxuan7/solution/pkg/utils"
	"time"
)

type Entity struct {
	utils.Model
	UID              string    `json:"uid"`
	PlatformUID      string    `json:"platform_uid"`
	CompanyId        uint      `json:"company_id"`
	ShopId           uint      `json:"shop_id"`
	PlatformView     string    `json:"platform_view"`
	PaidAt           time.Time `json:"paid_at"`
	Status           string    `json:"status"`
	ShopPlatformCode string    `json:"shop_platform_code"`
	WarehouseId      uint      `json:"warehouse_id"`
	CreatedMethod    string    `json:"created_method"`
}
