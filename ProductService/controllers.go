package main

import (
	"Assignment/OrderService/models"
	"github.com/gin-gonic/gin"
)

func GetCatalogue(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"Status": "Success",
		"Data":   models.Catalogue,
	})
}

func UpdateCatalogue(ctx *gin.Context) {
	var parsed map[string]interface{}
	err := ctx.BindJSON(&parsed)
	if err != nil {
		ctx.JSON(400, gin.H{
			"Status":  "ERROR",
			"Message": "something went wrong",
		})
	}
	for _, v := range parsed["products"].([]map[string]interface{}) {
		x := models.Catalogue[string(v["prodId"].(int64))]
		x.Availability--
		models.Catalogue[string(v["prodId"].(int64))] = x
	}
	ctx.JSON(200, models.Catalogue)
}
