package rule

import (
	"context"
	"github.com/wenxuan7/solution/pkg/utils"
)

const (
	StyleProduct   = "Product"   //商品匹配
	StyleFreebies  = "Freebies"  //赠品匹配
	StyleTag       = "Tag"       //标签匹配
	StyleException = "Exception" //异常匹配
	StyleWarehouse = "Warehouse" //仓库匹配
	StyleExpress   = "Express"   //快递匹配
	StyleMerge     = "Merge"     //合单规则匹配
	StyleSplit     = "Split"     //拆单规则匹配
)

type Reader interface {
}

type Writer interface {
	Set(ctx context.Context, r *Rule) error
}

type ReadWriter interface {
	Reader
	Writer
}

type Rule struct {
	utils.Model
	CompanyId        uint   `json:"company_id"`
	Style            string `json:"style"`
	ShopPlatformCode string `json:"shop_platform_code"`
	Name             string `json:"name"`
	Index            string `json:"index"`
	Opened           int    `json:"opened"`
}

func (r *Rule) TableName() string {
	return "rule"
}

type Multi struct {
	utils.Model
	RuleId uint   `json:"rule_id"`
	K      string `json:"k"`
	V      string `json:"v"`
}

func (m *Multi) TableName() string {
	return "rule_multi"
}
