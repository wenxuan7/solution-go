package settings

import (
	"context"
	"github.com/wenxuan7/solution/utils"
	"testing"
)

func TestService_SetGetSetsGets(t *testing.T) {
	e := &Entity{
		K: OrderSync,
		V: `{"sync_wait_payment_order":true}`,
	}
	ctx := utils.WithCompanyId(context.Background(), uint(1693))

	err := s.Set(ctx, e)
	if err != nil {
		t.Fatalf("settings: fail to Set in TestService_SetGetSetsGets: %v", err)
	}
	get, err := s.Get(ctx, OrderSync)
	if err != nil {
		t.Fatalf("settings: fail to Get in TestService_SetGetSetsGets: %v", err)
	}
	t.Logf("settings: success to Get in TestService_SetGetSetsGets: %v", get)

	e.V = `{"sync_wait_payment_order":false}`
	es := []*Entity{e}
	err = s.Sets(ctx, es)
	if err != nil {
		t.Fatalf("settings: fail to Sets in TestService_SetGetSetsGets: %v", err)
	}
	gets, err := s.Gets(ctx, []string{OrderSync})
	if err != nil {
		t.Fatalf("settings: fail to Gets in TestService_SetGetSetsGets: %v", err)
	}
	t.Logf("settings: success to Gets in TestService_SetGetSetsGets: %v", gets)
}
