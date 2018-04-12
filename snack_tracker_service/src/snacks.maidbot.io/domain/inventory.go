package domain

type InventoryChange struct {
  Quantity          int     `json:"quantity"`
  Direction         int     `json:"direction"`
  ItemCode          int     `json:"item_code"`
  CreatedAt         *int    `json:"created_at"`
}

type Item struct {
  Name            string  `json:"name"`
  Code            int     `json:"code"`
  CreatedAt       *int     `json:"created_at"`
  UpdatedAt       *int     `json:"updated_at"`
}
