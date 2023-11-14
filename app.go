package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// var DbName = "inventory"
// var DbUser = "root"
// var DbPassword = "golang123"

func (app *App) Initilise(DbUser string, DbPassword string, DbName string) error {
	ConnectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DbName)
	var err error
	app.DB, err = sql.Open("mysql", ConnectionString)
	if err != nil {
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()

	return nil

}

func (app *App) Run(address string) {
	http.ListenAndServe(address, app.Router)
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{Id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/put/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		// log.Println("Error decoding JSON:", err)
		SendError(w, http.StatusBadRequest, "Invalid Id")
		return
	}

	p := product{Id: key}
	err = p.deleteProduct(app.DB)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err.Error())
	}

	SendResponse(w, http.StatusOK, map[string]string{"result": "successful deletion"})

}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	log.Printf("End Point Hit get product")

	products, err := getProducts(app.DB)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendResponse(w, http.StatusOK, products)
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("End Point Hit get product with id ")
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["Id"])
	if err != nil {
		SendError(w, http.StatusBadRequest, "Not Found Product")
		return
	}

	p := product{Id: key}

	err = p.getProduct(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			SendError(w, http.StatusNotFound, "Product Not found")
		default:
			SendError(w, http.StatusInternalServerError, err.Error())

		}
		return
	}

	SendResponse(w, http.StatusOK, p)

}

func SendResponse(w http.ResponseWriter, StatusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(StatusCode)
	w.Write(response)

}

func SendError(w http.ResponseWriter, StatusCode int, err string) {
	error_messages := map[string]string{"error": err}
	SendResponse(w, StatusCode, error_messages)
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product

	//this is first insilitising the Decode method from
	//NewDecode function. It is copying the decoded data to the
	//struct p

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		// log.Println("Error decoding JSON:", err)
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = p.createProduct(app.DB)
	if err != nil {
		log.Println("Error creating product:", err)
		SendError(w, http.StatusCreated, err.Error())
		return
	}

}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		// log.Println("Error decoding JSON:", err)
		SendError(w, http.StatusBadRequest, "Invalid Id")
		return
	}

	var p product

	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		// log.Println("Error decoding JSON:", err)
		SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	p.Id = key
	err = p.updateProduct(app.DB)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err.Error())
		return

	}

	SendResponse(w, http.StatusOK, p)

}
