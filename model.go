package main

import (
	"database/sql"
	"fmt"
)

type product struct {
	Id       int     `json:"Id"`
	Name     string  `json:"Name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (p *product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("Delete from products where id =%v", p.Id)
	_, err := db.Exec(query)
	return err
}

func getProducts(db *sql.DB) ([]product, error) {
	Query := "SELECT id, name, price, quantity FROM Products"

	rows, err := db.Query(Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Defer closing the rows to ensure it happens after the function returns.

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.Id, &p.Name, &p.Price, &p.Quantity); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT name, quantity, price FROM products WHERE id = %d", p.Id)
	row := db.QueryRow(query)

	// Use pointers for the fields to be scanned
	var name string
	var price float64
	var quantity int

	err := row.Scan(&p.Id, &p.Name, &p.Quantity)
	if err != nil {
		return err
	}
	p.Name = name
	p.Price = price
	p.Quantity = quantity

	return nil
}

//creating function for POST method
// This function is being made with the struct

func (p *product) createProduct(db *sql.DB) error {
	fmt.Println("Create Product Hit")

	query := "INSERT INTO products (Id, Name, Quantity, Price) VALUES (?, ?, ?, ?)"

	result, err := db.Exec(query, p.Id, p.Name, p.Quantity, p.Price)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.Id = int(id)
	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("update products set name =%v, quantity=%v, price=%v where id = %v", p.Name, p.Quantity, p.Price, p.Id)

	//we have written the query know we need to execute it as well
	result, _ := db.Exec(query)

	effectedrows, err := result.RowsAffected()

	if effectedrows == 0 {
		fmt.Printf("The row doent exist!")
	}
	return err

}
