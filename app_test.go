package main

import (
	"fmt"
	"log"
	"net/http"
	"testing"

)

//we can define the maintest func one time only

var a App

func TestMain(m *testing.M){
	err := a.Initilise(DbUser, DbPassword, DbName)

	if err!=nil{
		log.Fatal("Error Occured while intilising the db")
	}

	m.Run()


}

func createTable(){
	createTableQuery := `Create Table If not Exists Products (
		id int NOT NULL auto_increment,
		name varchar(255) NOT NULL,
		quantity int,
		price float(10,7),
		PRIMARY KEY (id)
	);`

	_, err := a.DB.Exec(createTableQuery)

	if err!=nil{
		log.Fatal(err)
	}


}

func cleartable(){
	a.DB.Exec("Delete from products")
}

func addproduct(name string, quantity int, price float64){
	query := fmt.Sprintf("Insert into products(name, quantity, price) values('%v', '%v', '%v')", name, quantity, price)
	a.DB.Exec(query)
}
func TestGetProduct(t *testing.T){
	cleartable()
	addproduct("cup", 1, 100)
	http.NewRequest("Get", )

}
