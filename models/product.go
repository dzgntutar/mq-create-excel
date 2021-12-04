package models

type Product struct {
	Id    int32   `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
	Stock int32   `json:"stock"`
}
