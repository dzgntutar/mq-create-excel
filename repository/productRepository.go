package repository

import (
	"database/sql"
	"fmt"
	"mq-create-excel/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func (productRepository ProductRepository) Insert(product models.Product) error {
	fmt.Println(product)
	stmt, err := productRepository.DB.Prepare("insert into product (name,price,stock) values($1,$2,$3) ")
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(product.Name, product.Price, product.Stock)
	if err != nil {
		return err
	}

	return nil
}
