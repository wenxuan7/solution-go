package utils

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func WithCompanyId(ctx context.Context, companyId uint) context.Context {
	return context.WithValue(ctx, "companyId", companyId)
}

func GetCompanyId(ctx context.Context) uint {
	return ctx.Value("companyId").(uint)
}

func WithTraceId(ctx context.Context) (context.Context, error) {
	traceId, err := uuid.NewUUID()
	if err != nil {
		return ctx, fmt.Errorf("utils: fail to uuid.NewUUID in WithTraceId: %w", err)
	}
	return context.WithValue(ctx, "traceId", traceId.String()), nil
}

func GetTraceId(ctx context.Context) string {
	return ctx.Value("traceId").(string)
}
