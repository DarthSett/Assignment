package main

import (
	"Assignment/OrderService/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	models.Catalogue = map[string]models.Product{
		"1": {
			Id:           1,
			Availability: 10,
			Price:        20,
			Category:     "Premium",
		},
		"2": {
			Id:           2,
			Availability: 10,
			Price:        10,
			Category:     "Regular",
		},
		"3": {
			Id:           3,
			Availability: 10,
			Price:        5,
			Category:     "Budget",
		},
	}
	fmt.Printf("\n%+v", models.Catalogue)
	r := gin.Default()
	r.GET("/catalogue", GetCatalogue)
	r.POST("/update-catalogue", UpdateCatalogue)
	r.Run(":8081")
}

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
	fmt.Printf("%+v", parsed)
	for _, v := range parsed["products"].([]interface{}) {
		id := int64(v.(map[string]interface{})["ProdId"].(float64))
		fmt.Printf("\nid: %v", strconv.FormatInt(id, 10))
		x := models.Catalogue[strconv.FormatInt(id, 10)]
		fmt.Printf("\nx: %+v", x)
		x.Availability--
		models.Catalogue[strconv.FormatInt(id, 10)] = x
	}
	fmt.Printf("\n%+v", models.Catalogue)
	ctx.JSON(200, gin.H{"Data": models.Catalogue})
}
