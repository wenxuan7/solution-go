package rule

type Index interface {
	Match(target Index) bool
}

type DefaultIndex struct {
	ShopIds []uint `json:"shop_ids"`
}

type ProductIndex struct {
	DefaultIndex
}

type FreebiesIndex struct {
	DefaultIndex
}

type TagIndex struct {
	DefaultIndex
}

type ExceptionIndex struct {
	DefaultIndex
}

type WarehouseIndex struct {
	DefaultIndex
}

type ExpressIndex struct {
	DefaultIndex
}

type MergeIndex struct {
	DefaultIndex
}

type SplitIndex struct {
	DefaultIndex
}

type Definition struct {
	Key      string
	IndexImp Index
}

var (
	name =
)

var (
	Keys = map[string]*Definition{

	}
)

