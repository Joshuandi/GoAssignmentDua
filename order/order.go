package order

import "time"

type Order struct {
	Order_id      int
	Customer_name string
	Ordered_at    time.Time
}
