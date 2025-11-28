package utils

import (
	"context"
	"fmt"
)

func WithCompanyId(ctx context.Context, companyId uint) context.Context {
	return context.WithValue(ctx, "companyId", companyId)
}

func GetCompanyId(ctx context.Context) uint {
	return ctx.Value("companyId").(uint)
}

func WithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, "traceId", traceId)
}

func WithAutoTraceId(ctx context.Context) (context.Context, error) {
	traceId, err := UUID()
	if err != nil {
		return ctx, fmt.Errorf("utils: fail to UUID in WithAutoTraceId; %w", err)
	}
	return context.WithValue(ctx, "traceId", traceId), nil
}

func GetTraceId(ctx context.Context) string {
	return ctx.Value("traceId").(string)
}
