package domain

type InventoryChange struct {
  Quantity          int       `json:"quantity"`
  Mode              int       `json:"mode"`
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

type InventoryAggregate struct {
  ItemCode          string                   `json:"item_code"`
  Quantity          int                      `json:"quantity"`
  InventoryChanges  []*InventoryChange       `json:"inventory_changes"`
}
