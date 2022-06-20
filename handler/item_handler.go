package handler

import (
	"GoAssignmentDua/database"
	item "GoAssignmentDua/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ItemHandlerInterface interface {
	ItemHandler(w http.ResponseWriter, r *http.Request)
}

type ItemHandler struct {
	db    *sql.DB
	Items []item.Item
}

func NewItemHandler(db *sql.DB) ItemHandlerInterface {
	return &ItemHandler{db: db}
}

func (i *ItemHandler) ItemHandler(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	id := param["id"]

	switch r.Method {
	case http.MethodGet:
		if id != "" {
			i.getItemByIdHanlder(w, r, id)
		} else {
			i.getAllItemHanlder(w, r)
		}
	case http.MethodPost:
		i.postItemHandler(w, r)
	case http.MethodPut:
		i.updateItemHandler(w, r, id)
	case http.MethodDelete:
		i.deleteItemHandler(w, r, id)
	}
}

func (i *ItemHandler) getItemByIdHanlder(w http.ResponseWriter, r *http.Request, id string) {
	var results = []item.Item{}
	if id != "" {
		sqlSt := `Select * from items where id = $1`
		rows, err := database.Db.Query(sqlSt, id)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			var items = item.Item{}
			if err = rows.Scan(&items.Item_id,
				&items.Item_code,
				&items.Description,
				&items.Quantity,
				&items.Order_id); err != nil {
				fmt.Println("No Data", err)
			}
			results = append(results, items)
		}
		jsonData, _ := json.Marshal(&results)
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonData)
	}
}
func (i *ItemHandler) getAllItemHanlder(w http.ResponseWriter, r *http.Request) {}
func (i *ItemHandler) postItemHandler(w http.ResponseWriter, r *http.Request) {
	var items = item.Item{}
	json.NewDecoder(r.Body).Decode(&items)
	sqlSt := `Insert into items	
	(item_code, description, quantity)
	values ($1, $2, $3)
	Returning id ;`

	err := database.Db.QueryRow(sqlSt,
		items.Item_code,
		items.Description,
		items.Quantity,
	).Scan(&items.Item_id)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(fmt.Sprint("Item Created")))
	return
}
func (i *ItemHandler) updateItemHandler(w http.ResponseWriter, r *http.Request, id string) {}
func (i *ItemHandler) deleteItemHandler(w http.ResponseWriter, r *http.Request, id string) {}
