package models

import "time"

type Order struct {
	Id           int64          `json:"id"`
	OrderValue   float64        `json:"value"`
	DispatchDate time.Time      `json:"dispatchDate"`
	OrderStatus  string         `json:"status"`
	ProdQuantity int64          `json:"productQuantity"`
	Products     []OrderedProds `json:"products"`
}

type OrderedProds struct {
	ProdId int64
	Value  float64
	IsPrem bool
}

var Orders = make(map[string]Order)
