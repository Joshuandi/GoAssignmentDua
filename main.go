package main

import (
	"GoAssignmentDua/database"
	order_handler "GoAssignmentDua/handler"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var PORT = ":8088"

func main() {
	database.Db, database.Err = sql.Open("postgres", ConnectDbPsql(
		database.Host,
		database.User,
		database.Password,
		database.Dbname,
		database.Port))
	if database.Err != nil {
		panic(database.Err)
	}
	defer database.Db.Close()

	database.Err = database.Db.Ping()
	if database.Err != nil {
		panic(database.Err)
	}
	fmt.Println("Successfully Connect to Database")

	r := mux.NewRouter()
	orderHandler := order_handler.NewOrderHandler(database.Db)
	r.HandleFunc("/order", orderHandler.OrderHandler)
	r.HandleFunc("/order/{id}", orderHandler.OrderHandler)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8088",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func ConnectDbPsql(host, user, password, dbname string, port int) string {
	psqlInfo := fmt.Sprintf("host= %s port= %d user= %s "+
		" password= %s dbname= %s sslmode=disable",
		database.Host,
		database.Port,
		database.User,
		database.Password,
		database.Dbname)
	return psqlInfo
}
