package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT id, name, quantity, price from products"
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func getProductByID(db *sql.DB, id int) (*product, error) {
	query := "SELECT id, name, quantity, price FROM products WHERE id = $1"
	row := db.QueryRow(query, id)

	var p product
	err := row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *product) createProduct(db *sql.DB) error {
	query := "INSERT INTO products(name, quantity, price) VALUES($1, $2, $3) RETURNING id"

	err := db.QueryRow(query, p.Name, p.Quantity, p.Price).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE products set name='%v', quantity=%v, price=%v where id=%v", p.Name, p.Quantity, p.Price, p.ID)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such row exists")
	}

	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("DELETE from products where id=%v", p.ID)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such row exists")
	}

	return err
}
