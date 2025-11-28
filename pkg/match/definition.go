package match

import (
	"context"
	"github.com/wenxuan7/solution/pkg/order"
	"github.com/wenxuan7/solution/pkg/order/line"
)

type Matcher interface {
	Match(ctx context.Context, orderToLine map[*order.Entity][]*line.Entity) error
}
