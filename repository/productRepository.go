package repository

import (
	"database/sql"
	"mq-create-excel/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func (repository ProductRepository) Insert(product models.Product) error {
	stmt, err := repository.DB.Prepare("insert into product (name,price,stock) values($1,$2,$3) ")
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

func (repository ProductRepository) GetAll() ([]models.Product, error) {
	var product models.Product
	var products []models.Product

	rows, err := repository.DB.Query("select id , name, price , stock from product")

	if err != nil {
		return products, err
	}
	for rows.Next() {
		err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.Stock)
		if err != nil {
			return products, err
		}

		products = append(products, product)
	}

	return products, nil
}
