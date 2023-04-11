package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/catalogue", GetCatalogue)
	r.POST("/update-catalogue", UpdateCatalogue)
	r.Run("8080")
}
