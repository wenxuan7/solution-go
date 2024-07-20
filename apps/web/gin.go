package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wenxuan7/solution/utils"
	"log/slog"
	"net/http"
)

type resp gin.H

func newFailResp(traceId string) *resp {
	return &resp{
		"status":   "failure",
		"trace_id": traceId,
	}
}

func newSuccessResp(traceId string) *resp {
	return &resp{
		"status":   "success",
		"trace_id": traceId,
	}
}

func (r *resp) withMsg(s string) *resp {
	(*r)["msg"] = s
	return r
}

func (r *resp) withError(err string) *resp {
	(*r)["error"] = err
	return r
}

func (r *resp) withData(data any) *resp {
	(*r)["data"] = data
	return r
}

func getCtx(c *gin.Context) context.Context {
	ctx, err := utils.WithTraceId(context.Background())
	if err != nil {
		slog.Error("生成traceId错误", "error", err)
		c.JSON(http.StatusBadRequest, newFailResp("").withError("业务错误"))
		return nil
	}
	ctx = utils.WithCompanyId(ctx, 1693)
	return ctx
}
