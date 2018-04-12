package main

import (
	"encoding/json"
	"fmt"
	"time"
	// "log"
	// "strings"
	"net/http"
	// "html/template"
	"io/ioutil"
	"github.com/gorilla/mux"

	// "app/src/snacks.maidbot.io/clients"
	"app/src/snacks.maidbot.io/domain"
	// TODO: add service layer if things get more complicated....
	"app/src/snacks.maidbot.io/data"
)

var ASSETS_DIR = "/go/src/app/src/snacks.maidbot.io/apps/snacktrackerapi/assets/"

type SnackTrackerState struct {
	Mode *int								`json:"mode"`
	ItemCount *int 					`json:"item_count"`
	ItemCode *string				`json:"item_code"`
}

type SnackTrackerApiResources struct {
	inventoryData data.InventoryData
	stateTrackerState *SnackTrackerState
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
	now := time.Now()
	nanos := now.UnixNano()
	updatedBefore := nanos / 1000000

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

	if sr.stateTrackerState.ItemCode == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var inventoryChange = domain.InventoryChange{*sr.stateTrackerState.ItemCount, -1 * (*sr.stateTrackerState.Mode), *sr.stateTrackerState.ItemCode, nil}
	dbResult, dbErr := sr.inventoryData.CreateInventoryChange(&inventoryChange)
	if dbErr != nil {
		http.Error(w, dbErr.Error(), 500)
		return
	}

	// Marshal
	output, err := json.Marshal(dbResult)
	sr.stateTrackerState.ItemCode = nil
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (sr *SnackTrackerApiResources) listInventoryChanges(w http.ResponseWriter, r *http.Request) {
	updatedAfter := 0
	now := time.Now()
	nanos := now.UnixNano()
	updatedBefore := nanos / 1000000

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

		if state_component == "item_code" {
			sr.stateTrackerState.ItemCode = stateChangeState.ItemCode
			if sr.stateTrackerState.Mode == &domain.CHECKOUT_MODE {
				var inventoryChange = domain.InventoryChange{1, domain.CHECKOUT_MODE, *sr.stateTrackerState.ItemCode, nil}
				_, dbErr := sr.inventoryData.CreateInventoryChange(&inventoryChange)
				if dbErr != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			}
		}

		if state_component == "item_count" {
			sr.stateTrackerState.ItemCode = stateChangeState.ItemCode
		}

		if state_component == "mode" {
			oldMode := sr.stateTrackerState.Mode
			if *stateChangeState.Mode == 1 || *stateChangeState.Mode == -1 {
				sr.stateTrackerState.Mode = stateChangeState.Mode
				if oldMode != sr.stateTrackerState.Mode {
					sr.stateTrackerState.ItemCode = nil
				}
			}
		}

		// Marshal
		output, err := json.Marshal(sr.stateTrackerState)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
}

func (sr *SnackTrackerApiResources) getState(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Getting application state");

		// Marshal
		output, err := json.Marshal(sr.stateTrackerState)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
}


// func (sr *SnackTrackerApiResources) updateItem(w http.ResponseWriter, r *http.Request) {
// 	// TODO -- only parse once, this is to facilitate debugging.
// 	tmpl := template.Must(template.ParseFiles(ASSETS_DIR + "templates/software_status.html"))
// 	if r.Method != http.MethodPost {
// 		tmpl_err := tmpl.Execute(w, sr.software_status_data)
// 		if(tmpl_err != nil) {
// 			log.Print(tmpl_err)
// 		}
// 		return
// 	}
//
// 	resource_type := mux.Vars(r)["resource_type"]
// 	if resource_type == "app_status" {
// 		app_name := r.FormValue("app_name")
// 		// Hard code for testing
// 		app_status, app, app_err := sr.software_status_serivice.GetApplicationStatus(app_name)
// 		sr.software_status_data.Application = app
// 		sr.software_status_data.ApplicationStatus = app_status
// 		if(app_err != nil) {
// 			log.Print(app_err)
// 			sr.software_status_data.ApplicationStatusReady = false
// 		} else {
// 			sr.software_status_data.ApplicationStatusReady = true
// 		}
// 	}
//
// 	if resource_type == "robot_status" {
// 		robot_name := r.FormValue("robot_name")
// 		// Hard code for testing
// 		robot_status, robot, robot_err := sr.software_status_serivice.GetRobotStatus(robot_name)
// 		sr.software_status_data.Robot = robot
// 		sr.software_status_data.RobotStatus = robot_status
// 		if(robot_err != nil) {
// 			log.Print(robot_err)
// 			sr.software_status_data.RobotStatusReady = false
// 		} else {
// 			sr.software_status_data.RobotStatusReady = true
// 		}
// 	}
//
// 	tmpl_err := tmpl.Execute(w, sr.software_status_data)
// 	if(tmpl_err != nil) {
// 		log.Print(tmpl_err)
// 	}
// }

func hwHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

var snackTrackerState = SnackTrackerState{&domain.CHECKOUT_MODE, &domain.INTAKE_MODE, nil}

func main() {
  fmt.Printf("starting fleet snack tracker service \n")
	r := mux.NewRouter()
	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir(ASSETS_DIR)))
	inventoryData, _ := data.NewInventoryData()
	api_resources := &SnackTrackerApiResources{
		inventoryData,
		&snackTrackerState}

	r.HandleFunc("/", api_resources.landingPage)
	r.PathPrefix("/assets/").Handler(fs)
	r.HandleFunc("/ping", hwHandler)
	r.HandleFunc("/inventory_change", api_resources.createInventoryChange).Methods("POST")
	r.HandleFunc("/item", api_resources.createItem).Methods("POST")
	r.HandleFunc("/state/{state_component}", api_resources.setState).Methods("POST")
	r.HandleFunc("/state", api_resources.getState).Methods("GET")
	r.HandleFunc("/inventory_changes", api_resources.listInventoryChanges).Methods("GET")
	r.HandleFunc("/items", api_resources.listItems).Methods("GET")
	http.ListenAndServe(":80", r)
}
