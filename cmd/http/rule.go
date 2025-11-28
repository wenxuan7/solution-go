package main

import "github.com/gin-gonic/gin"

func ruleController(omsGroup *gin.RouterGroup) {
	group := omsGroup.Group("/rule", func(c *gin.Context) {})
	group.GET("/get_api_detail")
	group.POST("/get_api_detail")
}
