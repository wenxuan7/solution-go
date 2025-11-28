package main

import (
	"github.com/gin-gonic/gin"
	settings2 "github.com/wenxuan7/solution/pkg/settings"
	"github.com/wenxuan7/solution/pkg/utils"
	"log/slog"
	"net/http"
)

var settingsService *settings2.Service

func settingsController(omsGroup *gin.RouterGroup) {
	var err error
	settingsService, err = settings2.NewWithLCache()
	if err != nil {
		panic(err)
	}

	group := omsGroup.Group("/settings", func(c *gin.Context) {})
	group.GET("/get", settingsGet)
	group.POST("/gets", settingsGets)
	group.POST("/set", settingsSet)
	group.POST("/sets", settingsSets)
}

func settingsGet(c *gin.Context) {
	k := c.Query("k")
	if k == "" {
		c.JSON(http.StatusBadRequest, newFailResp("").withError("参数有误"))
		return
	}

	ctx := getCtx(c)
	if ctx == nil {
		return
	}

	get, err := settingsService.GetFromDb(ctx, k)
	if err != nil {
		slog.Error("查询配置错误", "error", err)
		c.JSON(http.StatusBadRequest, newFailResp(utils.GetTraceId(ctx)).withError("业务错误"))
		return
	}
	c.JSON(http.StatusOK, newSuccessResp(utils.GetTraceId(ctx)).withData(get))
}

func settingsGets(c *gin.Context) {
	ks := make([]string, 0)
	err := c.BindJSON(&ks)
	if err != nil {
		c.JSON(http.StatusBadRequest, newFailResp("").withError("参数错误"))
		return
	}

	ctx := getCtx(c)
	if ctx == nil {
		return
	}

	gets, err := settingsService.GetsFromDb(ctx, ks)
	if err != nil {
		slog.Error("批量查询配置错误", "error", err)
		c.JSON(http.StatusBadRequest, newFailResp(utils.GetTraceId(ctx)).withError("业务错误"))
		return
	}
	c.JSON(http.StatusOK, newSuccessResp(utils.GetTraceId(ctx)).withData(gets))
}

func settingsSet(c *gin.Context) {
	e := &settings2.Settings{}
	err := c.BindJSON(&e)
	if err != nil {
		c.JSON(http.StatusBadRequest, newFailResp("").withError("参数错误"))
		return
	}

	ctx := getCtx(c)
	if ctx == nil {
		return
	}

	err = settingsService.Set(ctx, e)
	if err != nil {
		slog.Error("保存配置错误", "error", err)
		c.JSON(http.StatusBadRequest, newFailResp(utils.GetTraceId(ctx)).withError("业务错误"))
		return
	}
	c.JSON(http.StatusOK, newSuccessResp(utils.GetTraceId(ctx)).withMsg("保存配置成功"))
}

func settingsSets(c *gin.Context) {
	es := make([]*settings2.Settings, 0)
	err := c.BindJSON(&es)
	if err != nil {
		c.JSON(http.StatusBadRequest, newFailResp("").withError("参数错误"))
		return
	}

	ctx := getCtx(c)
	if ctx == nil {
		return
	}

	err = settingsService.Sets(ctx, es)
	if err != nil {
		slog.Error("批量保存配置错误", "error", err)
		c.JSON(http.StatusBadRequest, newFailResp(utils.GetTraceId(ctx)).withError("业务错误"))
		return
	}
	c.JSON(http.StatusOK, newSuccessResp(utils.GetTraceId(ctx)).withMsg("批量保存配置成功"))
}
