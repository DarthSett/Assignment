package main

import (
	"Assignment/OrderService/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func GetOrder(ctx *gin.Context) {
	var order map[string]interface{}
	err := ctx.BindJSON(&order)
	if err != nil {
		ctx.JSON(400, gin.H{
			"Status":  "ERROR",
			"Message": "something went wrong",
		})
	}
	orderId := order["id"].(string)
	ctx.JSON(200, gin.H{
		"Status": "Success",
		"Data":   models.Orders[orderId],
	})
}

func PostOrder(ctx *gin.Context) {
	var parsed map[string]interface{}
	err := ctx.BindJSON(&parsed)
	if err != nil {
		ctx.JSON(400, gin.H{
			"Status":  "ERROR",
			"Message": "something went wrong",
		})
	}
	products := parsed["products"].([]models.OrderedProds)
	order := models.Order{
		Id:           int64(len(models.Orders)),
		OrderValue:   0,
		OrderStatus:  "Placed",
		ProdQuantity: 0,
		Products:     []models.OrderedProds{},
	}
	premCount := 0
	for _, v := range products {
		order.OrderValue = order.OrderValue + v.Value
		order.Products = append(order.Products, v)
		order.ProdQuantity++
		if v.IsPrem {
			premCount++
		}
	}
	if err := CheckCatalogue(order); err != nil {
		ctx.JSON(400, gin.H{
			"Status":  "ERROR",
			"Message": "something went wrong",
			"Error":   err,
		})
	}
	if premCount >= 3 {
		order.OrderValue = order.OrderValue - order.OrderValue/10
	}
	ctx.JSON(200, gin.H{
		"Status":  "Success",
		"Message": "Order Placed",
	})
}

func UpdateOrder(ctx *gin.Context) {
	var parsed map[string]interface{}
	err := ctx.BindJSON(&parsed)
	if err != nil {
		ctx.JSON(400, gin.H{
			"Status":  "ERROR",
			"Message": "something went wrong",
		})
	}
	orderId := parsed["id"].(int64)
	status := parsed["status"].(string)
	order := models.Orders[strconv.FormatInt(orderId, 10)]
	order.OrderStatus = status
	if status == "dispatched" {
		order.DispatchDate = time.Now()
	}
	models.Orders[strconv.FormatInt(orderId, 10)] = order
	ctx.JSON(200, gin.H{
		"Status":  "Success",
		"Message": "Order Updated",
	})
}

func CheckCatalogue(order models.Order) error {

	// Define the data to be sent as a map
	// Convert the data to JSON bytes
	jsonData, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", "https://productService/update-catalogue", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set the request header
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Print the response status code and body
	fmt.Println("Response status code:", resp.StatusCode)
	body := new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("Response body:", body.String())
	m := make(map[string]interface{})
	err = json.Unmarshal(body.Bytes(), &m)
	if err != nil {
		return err
	}
	for _, v := range order.Products {
		if m["Data"].(map[string]interface{})[strconv.FormatInt(v.ProdId, 10)].(int64) < 0 {
			return fmt.Errorf("not enough products left for productId %v", v.ProdId)
		}
	}
	return nil
}
