package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func main() {

	app := App{}
	app.Initilise(DbUser, DbPassword, DbName)
	app.Run("localhost:10000")

}
