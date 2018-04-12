package data

import (
	"database/sql"
	"time"
	"log"

	. "github.com/mattn/go-sqlite3"

	"app/src/snacks.maidbot.io/domain"
)

type InventoryData interface {
	CreateInventoryChange(inventoryChange *domain.InventoryChange) (*domain.InventoryChange, error)
	GetInventoryChangesByTime(createdAfter int, createdBefore int) ([]*domain.InventoryChange, error)
	DeleteInventoryChangesByTime(createdAfter int, createdBefore int) error

	CreateItem(item *domain.Item) (*domain.Item, error)
	UpdateItem(item *domain.Item) (*domain.Item, error)
	GetItem(code int) (*domain.Item, error)
	GetItem(name string) (*domain.Item, error)
	GetItemsByUpdatedTime(updatedAfter int, updatedBefore int) ([]*domain.Item, error)

	currentMillis() int
}

// This factory function is your "constructor" for your data layer.
func NewInventoryDataSqlite3(dataSourceName string) (InventoryData, error) {
	db, err := sql.Open("sqlite3", "/data/inventory.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS inventory_changes (created_at integer PRIMARY KEY, quantity integer NOT NULL, direction integer NOT NULL, item_code integer NOT NULL)")
	statement.Exec()
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS items (code integer PRIMARY KEY, name integer NOT NULL)")
	statement.Exec()

	if err != nil {
		return nil, err
	}

	return &inventoryData{
		db: db,
	}, nil
}

// progressData provides a non-working stub you can fill in.
// For a working example of progressData implemented in-memory, look at mock_item_data.go
type inventoryData struct {
	//placeholder for examples sake, we don't really care if you use an SQL db
	db *sql.DB
}

func (d *inventoryData) currentMillis() int {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}

func (d *inventoryData) CreateInventoryChange(inventoryChange *domain.InventoryChange) (*domain.InventoryChange, error) {
	statement, _ := db.Prepare("INSERT INTO inventory_changes (quantity, direction, item_code, created_at) VALUES (?, ?, ?, ?)")
	createdAt = d.currentMillis()
  result, err :=  statement.Exec(inventoryChange.Quantity, inventoryChange.Direction, inventoryChange.ItemCode, createdAt)
	if err != nil {
		return nil, err
	}

	inventoryChange.CreatedAt = createdAt
	return inventoryChange, nil
}

func (d *inventoryData) GetInventoryChangesByTime(createdAfter int, createdBefore int) ([]*domain.InventoryChange, error) {
	statement, _ := db.Prepare("SELECT quantity, direction, item_code, created_at FROM inventory_changes WHERE created_at>=? AND created_at<?")
	rows, err := statement.Exec(createdAfter, createdBefore)
	if err != nil {
		log.Println("Unable to get changes by date.")
		return nil, err
	}

	changeList := []*domain.InventoryChange{}
	for rows.Next() {
		inventoryChange := new(domain.InventoryChange)
		err = rows.Scan(&inventoryChange.Quantity, &inventoryChange.Direction, &inventoryChange.ItemCode, &inventoryChange.CreatedAt)
		changeList = append(changeList, inventoryChange)
	}

	return changeList, nil
}

func (d *inventoryData) DeleteInventoryChangesByTime(createdAfter int, createdBefore int) error {
	statement, _ := db.Prepare("DELETE FROM inventory_changes WHERE created_at>=? AND created_at<?")
	rows, err := statement.Exec(createdAfter, createdBefore)
	if err != nil {
		log.Println("Unable to get delete by date.")
		return err
	}

	return nil
}

func (d *inventoryData) CreateItem(item *domain.Item) (*domain.Item, error) {
	statement, _ := db.Prepare("INSERT INTO items (code, name, created_at, updated_at) VALUES (?, ?, ?, ?)")
	createdAt = d.currentMillis()
  result, err :=  statement.Exec(item.Code, item.Name, createdAt, createdAt)
	if err != nil {
		log.Println("Unable to create item.")
		return nil, err
	}

	item.CreatedAt = createdAt
	item.UpdatedAt = createdAt
	return item, nil
}

func (d *inventoryData) UpdateItem(item *domain.Item) (*domain.Item, error) {
	statement, _ := db.Prepare("UPDATE items SET name=?, updated_at=? WHERE code=?")
	updatedAt = d.currentMillis()
  result, err :=  statement.Exec(item.Name, updatedAt, item.Code)
	if err != nil {
		log.Println("Unable to update item.")
		return nil, err
	}

	item.UpdatedAt = updatedAt
	return item, nil
}

func (d *inventoryData) GetItem(code int) (*domain.Item, error) {
	var item domain.Item
	err := db.QueryRow("SELECT code, name, created_at, updated_at FROM items WHERE code=?", code).Scan(&item.Code, &item.Name, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		log.Println("Unable to get item by id.")
		return nil, err
	}

	return item, nil
}

func (d *inventoryData) GetItem(int code) (*domain.Item, error) {
	var item domain.Item
	err := db.QueryRow("SELECT code, name, created_at, updated_at FROM items WHERE Name=?", name).Scan(&item.Code, &item.Name, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		log.Println("Unable to get item by name.")
		return nil, err
	}

	return item, nil
}

func (d *inventoryData) GetItemsByUpdatedTime(createdAfter int, createdBefore int) ([]*domain.Item, error) {
	statement, _ := db.Prepare("SELECT code, name, created_at, updated_at FROM items WHERE updated_at>=? AND updated_at<?")
	rows, err := statement.Exec(createdAfter, createdBefore)
	if err != nil {
		log.Println("Unable to get items by date.")
		return nil, err
	}

	itemsList := []*domain.Item{}
	for rows.Next() {
		item := new(domain.Item)
		err = rows.Scan(&item.Name, &item.Code, &item.CreatedAt, &item.UpdatedAt)
		itemsList = append(itemsList, inventoryChange)
	}

	return itemsList, nil
}
