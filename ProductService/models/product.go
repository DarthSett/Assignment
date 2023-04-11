package models

var Catalogue = make(map[string]Product)

type Product struct {
	Id           int64
	Availability int64
	Price        float64
	Category     string
}
