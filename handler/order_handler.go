package handler

import (
	"GoAssignmentDua/database"
	item "GoAssignmentDua/model"
	order "GoAssignmentDua/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type OrderHandlerInterface interface {
	OrderHandler(w http.ResponseWriter, r *http.Request)
}

type OrderHandler struct {
	db *sql.DB
}

func NewOrderHandler(db *sql.DB) OrderHandlerInterface {
	return &OrderHandler{db: db}
}

func (o *OrderHandler) OrderHandler(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	id := param["id"]

	switch r.Method {
	case http.MethodGet:
		o.getAllOrderHanlder(w, r)
	case http.MethodPost:
		o.postOrderHandler(w, r)
	case http.MethodPut:
		o.updateOrderHandler(w, r, id)
	case http.MethodDelete:
		o.deleteOrderHandler(w, r, id)
	}
}
func (o *OrderHandler) getAllOrderHanlder(w http.ResponseWriter, r *http.Request) {
	sqlGet := `select
	o.order_id,
	o.customer_name,
	o.ordered_at,
	json_agg(json_build_object(
		'item_id',i.item_id,
		'item_code',i.item_code,
		'description',i.description,
		'quantity',i.quantity,
		'order_id',i.order_id
	)) as item
from public.orders o join items i on o.order_id = i.order_id
group by o.order_id;`
	rows, err := database.Db.Query(sqlGet)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var orders []*order.Order
	var items []item.Item

	for rows.Next() {
		var o order.Order
		var itemsStr string
		if err = rows.Scan(
			&o.Order_id,
			&o.Customer_name,
			&o.Ordered_at,
			&itemsStr,
		); err != nil {
			fmt.Println("No Data", err)
		}

		if err := json.Unmarshal([]byte(itemsStr), &items); err != nil {
			fmt.Errorf("Error when parsing data items")
		} else {
			o.Item = append(o.Item, items...)
		}
		orders = append(orders, &o)
	}
	jsonData, _ := json.Marshal(&orders)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)

}
func (o *OrderHandler) postOrderHandler(w http.ResponseWriter, r *http.Request) {
	var orders = order.Order{}
	json.NewDecoder(r.Body).Decode(&orders)
	sqlSt := `Insert into orders
	(customer_name, ordered_at)
	values ($1, $2)
	Returning order_id ;`

	err := database.Db.QueryRow(sqlSt,
		orders.Customer_name,
		time.Now(),
	).Scan(&orders.Order_id)
	if err != nil {
		panic(err)
	}
	fmt.Println(orders)
	//item
	for i := 0; i < len(orders.Item); i++ {
		orders.Item[i].Order_id = orders.Order_id
		sqlSt := `Insert into items
	(item_code, description, quantity, order_id)
	values ($1, $2, $3, $4)
	Returning item_id;`
		err := database.Db.QueryRow(sqlSt,
			orders.Item[i].Item_code,
			orders.Item[i].Description,
			orders.Item[i].Quantity,
			orders.Item[i].Order_id,
		).Scan(&orders.Item[i].Item_id)
		if err != nil {
			panic(err)
		}
	}
	w.Write([]byte(fmt.Sprint("Order Created")))
	return
}

func (o *OrderHandler) updateOrderHandler(w http.ResponseWriter, r *http.Request, id string) {
	for id != "" {
		var orders order.Order
		//var orders = order.Order{}
		json.NewDecoder(r.Body).Decode(&orders)
		sqlSt := `update orders set customer_name = $2, ordered_at = $3 where order_id = $1;`
		res, err := database.Db.Exec(sqlSt,
			id,
			orders.Customer_name,
			time.Now(),
		)
		if err != nil {
			panic(err)
		}
		fmt.Println("ini data order", orders)
		for i := 0; i < len(orders.Item); i++ {
			var items item.Item
			//orders.Item[i].Order_id = orders.Order_id
			items.Item_id = orders.Item[i].Item_id
			items.Item_code = orders.Item[i].Item_code
			items.Description = orders.Item[i].Description
			items.Quantity = orders.Item[i].Quantity
			fmt.Println("ini data items", items)
			sqlSts := `update items set item_code = $3, description = $4, quantity = $5
			where order_id = $1 and item_id = $2;`
			_, err := database.Db.Exec(sqlSts,
				id,
				items.Item_id,
				items.Item_code,
				items.Description,
				items.Quantity,
			)
			if err != nil {
				panic(err)
			}
		}
		count, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}
		w.Write([]byte(fmt.Sprint("Update data :", count)))
		return
	}
}
func (o *OrderHandler) deleteOrderHandler(w http.ResponseWriter, r *http.Request, id string) {
	sqlDelete := `DELETE from orders WHERE order_id = $1`
	if index, err := strconv.Atoi(id); err == nil {
		res, err := database.Db.Exec(sqlDelete, index)
		if err != nil {
			panic(err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}
		w.Write([]byte(fmt.Sprint("Deleted Data", count)))
		return
	}
}
