package main

import (
	"encoding/json"
	"fmt"
	"log"
	// "strings"
	"net/http"
	"html/template"
	"io/ioutil"
	"github.com/gorilla/mux"

	// "app/src/snacks.maidbot.io/clients"
	"app/src/snacks.maidbot.io/domain"
	// TODO: add service layer if things get more complicated....
	"app/src/snacks.maidbot.io/data"
)

var ASSETS_DIR = "/go/src/app/src/snacks.maidbot.io/apps/snacktrackerapi/assets/"

type SnackTrackerApiResources struct {
	inventory_data data.InventoryData
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
		dbResult, dbErr = sr.CreateInventoryChange(inventoryChange)
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
		dbResult, dbErr = sr.CreateInventoryChange(item)
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

// func (sr *ServiceResources) updateItem(w http.ResponseWriter, r *http.Request) {
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

func main() {
  fmt.Printf("starting fleet status service \n")
	r := mux.NewRouter()
	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir(ASSETS_DIR)))
	api_resources := &SnackTrackerApiResources{
		data.NewInventoryDataSqlite3()}

	r.PathPrefix("/assets/").Handler(fs)
	r.HandleFunc("/inventory_change", api_resources.createInventoryChange).Methods("POST")
	r.HandleFunc("/item", api_resources.createItem).Methods("POST")

	http.ListenAndServe(":80", r)
}
