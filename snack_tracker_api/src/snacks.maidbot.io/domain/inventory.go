package domain

type InventoryChange struct {
  Quantity          int       `json:"quantity"`
  Direction         int       `json:"direction"`
  ItemCode          string    `json:"item_code"`
  ItemName          *string   `json:"item_name"`
  CreatedAt         *int64    `json:"created_at"`
}

type Item struct {
  Code            string      `json:"code"`
  Name            string      `json:"name"`
  CreatedAt       *int64      `json:"created_at"`
  UpdatedAt       *int64      `json:"updated_at"`
}
