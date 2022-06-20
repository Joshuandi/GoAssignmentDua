package model

type Item struct {
	Item_id     int    `json:"item_id"`
	Item_code   string `json:"item_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Order_id    int    `json:"order_id"`
}
