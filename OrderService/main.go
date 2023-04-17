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

func main() {
	r := gin.Default()
	r.POST("/OrderGet", GetOrder)
	r.POST("/OrderPlace", PostOrder)
	r.POST("/OrderModify", UpdateOrder)
	r.Run(":8080")
}

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
	products := parsed["products"].([]interface{})
	order := models.Order{
		Id:           int64(len(models.Orders)) + 1,
		OrderValue:   0,
		OrderStatus:  "Placed",
		ProdQuantity: 0,
		Products:     []models.OrderedProds{},
	}
	premCount := 0
	for _, v := range products {
		order.OrderValue = order.OrderValue + v.(map[string]interface{})["Value"].(float64)
		orderedProd := models.OrderedProds{
			ProdId: int64(v.(map[string]interface{})["ProdId"].(float64)),
			Value:  v.(map[string]interface{})["Value"].(float64),
			IsPrem: v.(map[string]interface{})["IsPrem"].(bool),
		}
		order.Products = append(order.Products, orderedProd)
		order.ProdQuantity++
		if v.(map[string]interface{})["IsPrem"].(bool) {
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
	models.Orders[strconv.FormatInt(order.Id, 10)] = order
	fmt.Printf("\nORDERS LIST: %+v", models.Orders)
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
	orderId := int64(parsed["id"].(float64))
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
	bodyReq := make(map[string]interface{})
	bodyReq["products"] = order.Products
	jsonData, err := json.Marshal(bodyReq)
	println(string(jsonData))
	if err != nil {
		return err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", "http://127.0.0.1:8081/update-catalogue", bytes.NewBuffer(jsonData))
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
	fmt.Printf("%+v", m)
	for _, v := range order.Products {
		prod := m["Data"].(map[string]interface{})[strconv.FormatInt(v.ProdId, 10)].(map[string]interface{})
		if int64(prod["Availability"].(float64)) < 0 {
			return fmt.Errorf("not enough products left for productId %v", v.ProdId)
		}
	}
	return nil
}
