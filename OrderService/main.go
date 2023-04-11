package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.POST("/OrderGet", GetOrder)
	r.POST("/OrderPlace", PostOrder)
	r.POST("/OrderModify", UpdateOrder)
	r.Run("8080")
}
