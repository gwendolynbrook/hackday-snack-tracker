package main

import (
	"encoding/json"
	"fmt"
	"time"
	"log"
	"strconv"
	"os"
	// "strings"
	"net/http"
	"html/template"
	"io/ioutil"
	"github.com/gorilla/mux"

	"app/src/snacks.maidbot.io/domain"
	// TODO: add service layer if things get more complicated....
	"app/src/snacks.maidbot.io/data"
)

var ASSETS_DIR = "/go/src/app/src/snacks.maidbot.io/apps/snacktrackerapi/assets/"
var CACHE_DIR = "/go/src/app/src/snacks.maidbot.io/apps/snacktrackerapi/cache/"
var PLEASE_SCAN_SNACK = "Please Scan a Snack"

type SnackTrackerState struct {
	Mode int								`json:"mode"`
	ItemCount *int 					`json:"item_count"`
	ItemCode string					`json:"item_code"`
	ItemName *string 				`json:"item_name"`
	RemainingQuantity *int	`json:"remaining_quantity"`
	CodeIsNew bool					`json:"code_is_new"`
}

type SnackTrackerApiResources struct {
	inventoryData 					data.InventoryData
	snackTrackerState 			*SnackTrackerState
}

func currentMillis() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}

func (d *SnackTrackerState) cacheExists() bool {
	if _, err := os.Stat(CACHE_DIR); os.IsNotExist(err) {
		log.Print("No snack tracker cache directory; creating one!")
		os.MkdirAll(CACHE_DIR, 0755)
		return false
	}

	if _, err := os.Stat(CACHE_DIR + "snack_tracker_state.json"); os.IsNotExist(err) {
		log.Print("No snack tracker cache json!")
		return false
	}

	return true
}

func (d *SnackTrackerState) load() {
	if !(d.cacheExists()) {
		log.Print("No cache. Using default.")
		return
	}

	cachedState, readErr := ioutil.ReadFile(CACHE_DIR + "snack_tracker_state.json")
	if readErr != nil {
		log.Print("Cannot read cache; using defaults!")
		return
	}

	jsonErr := json.Unmarshal(cachedState, d)
	if jsonErr != nil {
		log.Print("Failed to unmarshal cache json.")
	}
}

func (d *SnackTrackerState) save() {
	stateToCache, jsonErr := json.Marshal(d)
	if jsonErr != nil {
		log.Print("Failed to mashal state json.")
		return
	}

  err := ioutil.WriteFile(CACHE_DIR + "snack_tracker_state.json", stateToCache, 0644)
	if err != nil {
		log.Print("Failed to write state to cache")
	}
}

func (sr *SnackTrackerApiResources) createInventoryChange(w http.ResponseWriter, r *http.Request) {
		// Read
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Unmarshal
		var inventoryChange domain.InventoryChange
		err = json.Unmarshal(b, &inventoryChange)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Update DB
		dbResult, dbErr := sr.inventoryData.CreateInventoryChange(&inventoryChange)
		if dbErr != nil {
			http.Error(w, dbErr.Error(), 500)
			return
		}

		// Marshal
		output, err := json.Marshal(dbResult)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
}

func (sr *SnackTrackerApiResources) listItems(w http.ResponseWriter, r *http.Request) {
	updatedAfter := 0
	updatedBefore := currentMillis()

	items, dbErr := sr.inventoryData.GetItemsByUpdatedTime(int64(updatedAfter), updatedBefore)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), 500)
		return
	}

	// Marshal
	output, err := json.Marshal(items)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (sr *SnackTrackerApiResources) undoLastInventoryChange(w http.ResponseWriter, r *http.Request) {

	if sr.snackTrackerState.ItemCode == PLEASE_SCAN_SNACK {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	log.Printf("Putting back the last inventory change for %s : %s", sr.snackTrackerState.ItemCode, *sr.snackTrackerState.ItemName)

	var inventoryChange = domain.InventoryChange{-1 * (*sr.snackTrackerState.ItemCount), sr.snackTrackerState.Mode, sr.snackTrackerState.ItemCode, nil, nil}
	dbResult, dbErr := sr.inventoryData.CreateInventoryChange(&inventoryChange)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), 500)
		return
	}

	// Marshal
	output, err := json.Marshal(dbResult)
	sr.snackTrackerState.ItemCode = PLEASE_SCAN_SNACK
	var remaining = *sr.snackTrackerState.RemainingQuantity + -1 * (*sr.snackTrackerState.ItemCount) * sr.snackTrackerState.Mode
	sr.snackTrackerState.RemainingQuantity = &remaining

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (sr *SnackTrackerApiResources) listInventoryChanges(w http.ResponseWriter, r *http.Request) {
	updatedAfter := 0
	updatedBefore := currentMillis()

	items, dbErr := sr.inventoryData.GetInventoryChangesByTime(int64(updatedAfter), updatedBefore)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), 500)
		return
	}

	// Marshal
	output, err := json.Marshal(items)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (sr *SnackTrackerApiResources) createItem(w http.ResponseWriter, r *http.Request) {
		// Read
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Unmarshal
		var item domain.Item
		err = json.Unmarshal(b, &item)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Update DB
		dbResult, dbErr := sr.inventoryData.CreateItem(&item)
		if dbErr != nil {
			http.Error(w, dbErr.Error(), 500)
			return
		}

		// Marshal
		output, err := json.Marshal(dbResult)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
}

func (sr *SnackTrackerApiResources) setState(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Setting application state");
		state_component := mux.Vars(r)["state_component"]

		// Read
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Unmarshal
		var stateChangeState SnackTrackerState
		err = json.Unmarshal(b, &stateChangeState)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// To bootstrap single-div refreshing via javascript
		if state_component == "code_is_new" {
			sr.snackTrackerState.CodeIsNew = false
		}

		// HACK HACK HACK -- the meaty bits; interaction with the barcode scanner
		if state_component == "item_code" {
			sr.snackTrackerState.ItemCode = stateChangeState.ItemCode
			if sr.snackTrackerState.Mode == domain.CHECKOUT_MODE && sr.snackTrackerState.ItemCode != PLEASE_SCAN_SNACK {
				log.Printf("Setting new item_code in CHECKOUT mode %s", sr.snackTrackerState.ItemCode)
				var inventoryChange = domain.InventoryChange{1, domain.CHECKOUT_MODE, sr.snackTrackerState.ItemCode, nil, nil}
				_, dbErr := sr.inventoryData.CreateInventoryChange(&inventoryChange)
				sr.snackTrackerState.ItemName = inventoryChange.ItemName
				if dbErr != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			} else {
				log.Print("Setting new item_code in INTAKE mode")
				sr.snackTrackerState.CodeIsNew = true
				item, dbErr := sr.inventoryData.GetItemByCode(sr.snackTrackerState.ItemCode)
				if dbErr == nil {
					sr.snackTrackerState.ItemName = &item.Name
				} else {
					sr.snackTrackerState.ItemName = nil
				}
			}

			// Compute the aggregate once we get the input
			inventoryAggregate, dbErr := sr.inventoryData.ComputeInventoryAggregate(sr.snackTrackerState.ItemCode, int64(0), currentMillis())
			if dbErr == nil {
				sr.snackTrackerState.RemainingQuantity = &inventoryAggregate.Quantity
			}
		}

		// Marshal
		output, err := json.Marshal(sr.snackTrackerState)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
}

func (sr *SnackTrackerApiResources) getState(w http.ResponseWriter, r *http.Request) {
		// Marshal
		output, err := json.Marshal(sr.snackTrackerState)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
}


func (sr *SnackTrackerApiResources) landingPageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO -- only parse once, this is to facilitate debugging.
	tmpl := template.Must(template.ParseFiles(ASSETS_DIR + "templates/landing.html"))
	// TODO! Template the redirects here!
	if r.Method != http.MethodPost {
		tmpl_err := tmpl.Execute(w, sr.snackTrackerState)
		if(tmpl_err != nil) {
			log.Print(tmpl_err)
		}
		return
	}
}

func (sr *SnackTrackerApiResources) addSnackInventoryHandler(w http.ResponseWriter, r *http.Request) {
	// TODO -- only parse once, this is to facilitate debugging.
	var zero = 0
	sr.snackTrackerState.Mode = domain.INTAKE_MODE
	sr.snackTrackerState.ItemCount = &zero
	tmpl := template.Must(template.ParseFiles(ASSETS_DIR + "templates/add_snack_inventory.html"))
	// TODO! Template the redirects here!
	if r.Method != http.MethodPost {
		sr.snackTrackerState.ItemCode = PLEASE_SCAN_SNACK
		sr.snackTrackerState.ItemName = nil
		sr.snackTrackerState.RemainingQuantity = nil
		tmpl_err := tmpl.Execute(w, sr.snackTrackerState)
		if(tmpl_err != nil) {
			log.Print(tmpl_err)
		}
		sr.snackTrackerState.save()
		return
	}

	r.ParseForm()
	log.Print(r.Form)

	item_code := sr.snackTrackerState.ItemCode
	log.Print("Processig intake item_code : " + item_code + ", name : " + r.FormValue("item_name") + ", count : ", r.FormValue("item_count"))


	if item_code == PLEASE_SCAN_SNACK || item_code == "" {
		log.Print("Cannot accept nil item code.")
		tmpl_err := tmpl.Execute(w, sr.snackTrackerState)
		if(tmpl_err != nil) {
			log.Print(tmpl_err)
		}
		return
	}

	item_name := r.FormValue("item_name")
	item_count, _ := strconv.Atoi(r.FormValue("item_count"))

	var item = domain.Item{item_code, item_name, nil, nil}
	if sr.snackTrackerState.ItemName == nil {
		_, err := sr.inventoryData.CreateItem(&item)
		if err != nil {
			fmt.Println("Cannot update item; will create in data layer")
		}
	} else if *sr.snackTrackerState.ItemName != item.Name {
		_, err := sr.inventoryData.UpdateItem(&item)
		if err != nil {
			fmt.Println("Cannot update item; will create in data layer")
		}
	}

	var inventoryChange = domain.InventoryChange{item_count, domain.INTAKE_MODE, item_code, &item_name, nil}
	stampedInventoryChange, err := sr.inventoryData.CreateInventoryChange(&inventoryChange)

	if err == nil {
		sr.snackTrackerState.ItemCode = stampedInventoryChange.ItemCode
		sr.snackTrackerState.ItemName = stampedInventoryChange.ItemName
		sr.snackTrackerState.ItemCount = &stampedInventoryChange.Quantity
		if sr.snackTrackerState.RemainingQuantity == nil {
			sr.snackTrackerState.RemainingQuantity = sr.snackTrackerState.ItemCount
		} else {
			var remaining = *sr.snackTrackerState.ItemCount + *sr.snackTrackerState.RemainingQuantity
			sr.snackTrackerState.RemainingQuantity = &remaining
		}
		sr.snackTrackerState.CodeIsNew = false
	} else {
		log.Print(err)
	}

	tmpl_err := tmpl.Execute(w, sr.snackTrackerState)
	if(tmpl_err != nil) {
		log.Print(tmpl_err)
	}
	return
}

func (sr *SnackTrackerApiResources) consumeSnacksHandler(w http.ResponseWriter, r *http.Request) {
	// TODO -- only parse once, this is to facilitate debugging.
	var one = 1
	sr.snackTrackerState.Mode = domain.CHECKOUT_MODE
	sr.snackTrackerState.ItemCount = &one
	tmpl := template.Must(template.ParseFiles(ASSETS_DIR + "templates/consume_snacks.html"))
	// TODO! Template the redirects here!
	if r.Method != http.MethodPost {
		sr.snackTrackerState.ItemCode = PLEASE_SCAN_SNACK
		sr.snackTrackerState.ItemName = nil
		sr.snackTrackerState.RemainingQuantity = nil
		tmpl_err := tmpl.Execute(w, sr.snackTrackerState)
		if(tmpl_err != nil) {
			log.Print(tmpl_err)
		}
		sr.snackTrackerState.save()
		return
	}
}

func hwHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

var snackTrackerState = SnackTrackerState{domain.CHECKOUT_MODE, nil, PLEASE_SCAN_SNACK, nil, nil, false}

func main() {
  fmt.Printf("starting fleet snack tracker service \n")
	r := mux.NewRouter()
	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir(ASSETS_DIR)))
	inventoryData, _ := data.NewInventoryData()
	snackTrackerState.load()
	api_resources := &SnackTrackerApiResources{
		inventoryData,
		&snackTrackerState}

	r.PathPrefix("/assets/").Handler(fs)
	r.HandleFunc("/", api_resources.landingPageHandler)
	r.HandleFunc("/consume_snacks", api_resources.consumeSnacksHandler)
	r.HandleFunc("/add_snack_inventory", api_resources.addSnackInventoryHandler)

	r.HandleFunc("/ping", hwHandler)
	r.HandleFunc("/inventory_change", api_resources.createInventoryChange).Methods("POST")
	r.HandleFunc("/inventory_change/undo", api_resources.undoLastInventoryChange).Methods("POST")
	r.HandleFunc("/item", api_resources.createItem).Methods("POST")
	r.HandleFunc("/state/{state_component}", api_resources.setState).Methods("POST")
	r.HandleFunc("/state", api_resources.getState).Methods("GET")
	r.HandleFunc("/inventory_changes", api_resources.listInventoryChanges).Methods("GET")
	r.HandleFunc("/items", api_resources.listItems).Methods("GET")
	http.ListenAndServe(":80", r)
}
