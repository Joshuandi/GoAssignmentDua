package model

import "time"

type Order struct {
	Order_id      int       `json:"order_id"`
	Customer_name string    `json:"customer_name"`
	Ordered_at    time.Time `json:"ordered_at"`
	Item          []Item    `json:"item"`
}
