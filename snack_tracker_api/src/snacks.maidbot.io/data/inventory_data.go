package data

import (
	"database/sql"
	"time"
	"log"
	"os"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"app/src/snacks.maidbot.io/domain"
)

type InventoryData interface {
	CreateInventoryChange(inventoryChange *domain.InventoryChange) (*domain.InventoryChange, error)
	GetInventoryChangesByTime(createdAfter int64, createdBefore int64) ([]*domain.InventoryChange, error)
	DeleteInventoryChangesByTime(createdAfter int64, createdBefore int64) error

	CreateItem(item *domain.Item) (*domain.Item, error)
	UpdateItem(item *domain.Item) (*domain.Item, error)
	GetItemByCode(code string) (*domain.Item, error)
	GetItemByName(name string) (*domain.Item, error)
	GetItemsByUpdatedTime(updatedAfter int64, updatedBefore int64) ([]*domain.Item, error)

	currentMillis() int64
}

// This factory function is your "constructor" for your data layer.
func NewInventoryData() (InventoryData, error) {
	data_dir := os.Getenv(DATA_DIR_ENV)
	os.MkdirAll(data_dir, 0755)
  os.Create(data_dir + "/inventory.db")
	fmt.Println("Using data dir : " + data_dir)
	db, err := sql.Open("sqlite3", data_dir + "/inventory.db")
	log.Print(err)
	statement, stmtErr := db.Prepare("CREATE TABLE IF NOT EXISTS inventory_changes (created_at integer PRIMARY KEY, quantity integer NOT NULL, direction integer NOT NULL, item_code text NOT NULL, item_name text NOT NULL)")
	log.Print(stmtErr)
	statement.Exec()
	statement, stmtErr = db.Prepare("CREATE TABLE IF NOT EXISTS items (code text PRIMARY KEY, name text NOT NULL, created_at integer NOT NULL, updated_at integer NOT_NULL)")
	log.Print(stmtErr)
	statement.Exec()
	db.Close()
	db, err = sql.Open("sqlite3", data_dir + "/inventory.db")

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

func (d *inventoryData) currentMillis() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}

func (d *inventoryData) CreateInventoryChange(inventoryChange *domain.InventoryChange) (*domain.InventoryChange, error) {
	item, itemErr := d.GetItemByCode(inventoryChange.ItemCode)
	if itemErr != nil {
		var item = domain.Item{inventoryChange.ItemCode, UNREGISTERED_NAME, nil, nil}
		newItem, newItemErr := d.CreateItem(&item)
		if newItemErr != nil {
			log.Println("Error creating item for inventory change")
			return nil, newItemErr
		}
		inventoryChange.ItemName = &newItem.Name
	} else {
		inventoryChange.ItemName = &item.Name
	}

	statement, _ := d.db.Prepare("INSERT INTO inventory_changes (quantity, direction, item_code, item_name, created_at) VALUES (?, ?, ?, ?, ?)")
	createdAt := d.currentMillis()
  _, err :=  statement.Exec(inventoryChange.Quantity, inventoryChange.Direction, inventoryChange.ItemCode, inventoryChange.ItemName, createdAt)
	if err != nil {
		return nil, err
	}

	inventoryChange.CreatedAt = &createdAt
	return inventoryChange, nil
}

func (d *inventoryData) GetInventoryChangesByTime(createdAfter int64, createdBefore int64) ([]*domain.InventoryChange, error) {
	rows, err := d.db.Query("SELECT quantity, direction, item_code, item_name, created_at FROM inventory_changes WHERE created_at>=? AND created_at<?", createdAfter, createdBefore)
	if err != nil {
		log.Println("Unable to get changes by date.")
		return nil, err
	}

	changeList := []*domain.InventoryChange{}
	for rows.Next() {
		inventoryChange := new(domain.InventoryChange)
		err = rows.Scan(&inventoryChange.Quantity, &inventoryChange.Direction, &inventoryChange.ItemCode, &inventoryChange.ItemName, &inventoryChange.CreatedAt)
		changeList = append(changeList, inventoryChange)
	}

	return changeList, nil
}

func (d *inventoryData) DeleteInventoryChangesByTime(createdAfter int64, createdBefore int64) error {
	statement, err := d.db.Prepare("DELETE FROM inventory_changes WHERE created_at>=? AND created_at<?")
	_, err = statement.Exec(createdAfter, createdBefore)
	if err != nil {
		log.Println("Unable to get delete by date.")
		return err
	}

	return nil
}

func (d *inventoryData) CreateItem(item *domain.Item) (*domain.Item, error) {
	statement, _ := d.db.Prepare("INSERT INTO items (code, name, created_at, updated_at) VALUES (?, ?, ?, ?)")
	createdAt := d.currentMillis()
  _, err :=  statement.Exec(item.Code, item.Name, createdAt, createdAt)
	if err != nil {
		log.Println("Unable to create item.")
		return nil, err
	}

	item.CreatedAt = &createdAt
	item.UpdatedAt = &createdAt
	return item, nil
}

func (d *inventoryData) UpdateItem(item *domain.Item) (*domain.Item, error) {
	statement, _ := d.db.Prepare("UPDATE items SET name=?, updated_at=? WHERE code=?")
	updatedAt := d.currentMillis()
  _, err :=  statement.Exec(item.Name, updatedAt, item.Code)
	if err != nil {
		log.Println("Unable to update item.")
		return nil, err
	}

	item.UpdatedAt = &updatedAt
	return item, nil
}

func (d *inventoryData) GetItemByCode(code string) (*domain.Item, error) {
	var item domain.Item
	err := d.db.QueryRow("SELECT code, name, created_at, updated_at FROM items WHERE code=?", code).Scan(&item.Code, &item.Name, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		log.Println("Unable to get item by code.")
		return nil, err
	}

	return &item, nil
}

func (d *inventoryData) GetItemByName(name string) (*domain.Item, error) {
	var item domain.Item
	err := d.db.QueryRow("SELECT code, name, created_at, updated_at FROM items WHERE code=?", name).Scan(&item.Code, &item.Name, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		log.Println("Unable to get item by name.")
		return nil, err
	}
	return &item, nil
}

func (d *inventoryData) GetItemsByUpdatedTime(createdAfter int64, createdBefore int64) ([]*domain.Item, error) {
	rows, err := d.db.Query("SELECT code, name, created_at, updated_at FROM items WHERE updated_at>=? AND updated_at<?", createdAfter, createdBefore)
	if err != nil {
		log.Println("Unable to get items by date.")
		return nil, err
	}

	itemsList := []*domain.Item{}
	for rows.Next() {
		item := new(domain.Item)
		err = rows.Scan(&item.Code, &item.Name, &item.CreatedAt, &item.UpdatedAt)
		itemsList = append(itemsList, item)
	}

	return itemsList, nil
}
