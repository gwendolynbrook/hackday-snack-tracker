package domain

type InventoryChange struct {
  Quantity          int       `json:"quantity"`
  Direction         int       `json:"direction"`
  ItemCode          string    `json:"item_code"`
  CreatedAt         *int64    `json:"created_at"`
}

type Item struct {
  Name            string      `json:"name"`
  Code            string      `json:"code"`
  CreatedAt       *int64      `json:"created_at"`
  UpdatedAt       *int64      `json:"updated_at"`
}
