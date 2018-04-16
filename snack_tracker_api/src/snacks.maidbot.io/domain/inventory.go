package domain

import (
  "strconv"
)

type InventoryChange struct {
  Quantity          int       `json:"quantity"`
  Mode              int       `json:"mode"`
  ItemCode          string    `json:"item_code"`
  ItemName          *string   `json:"item_name"`
  CreatedAt         *int64    `json:"created_at"`
}

func (c *InventoryChange) GetHeaders() []string {
  return []string{"quantity", "mode", "item_code", "item_name", "create_at"}
}

func (c *InventoryChange) ToSlice() []string {
  return []string{
    strconv.Itoa(c.Quantity),
    strconv.Itoa(c.Mode),
    c.ItemCode,
    *c.ItemName,
    strconv.Itoa(int(*c.CreatedAt))}
}

type Item struct {
  Code            string      `json:"code"`
  Name            string      `json:"name"`
  CreatedAt       *int64      `json:"created_at"`
  UpdatedAt       *int64      `json:"updated_at"`
}

type InventoryAggregate struct {
  ItemCode          string                   `json:"item_code"`
  ItemName          string                   `json:"item_name"`
  Quantity          int                      `json:"quantity"`
  InventoryChanges  []*InventoryChange       `json:"inventory_changes"`
}

func (c *InventoryAggregate) GetHeaders() []string {
  return []string{"quantity", "item_code", "item_name"}
}

func (c *InventoryAggregate) ToSlice() []string {
  return []string{
    strconv.Itoa(c.Quantity),
    c.ItemCode,
    c.ItemName}
}
