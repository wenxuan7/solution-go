package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func StartGin() {
	var (
		err error
	)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err = r.Run(); err != nil {
		panic(err)
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
