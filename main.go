package main

import (
	"encoding/json"
	"fmt"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "study/dom/docs"

	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "TODO",
		Description: "Implement CRUD API",
	},
}
// createEvent godoc
// @Summary create details of event
// @Description create details of event
// @Tags event
// @Success 200 {object} event
// @Router /event [post]
func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Schema wrong!")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

// getOneEvent godoc
// @Summary Get details of one event
// @Description Get details of one event
// @Tags event
// @Success 200 {object} event
// @Router /events/{id} [get]
func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

// getAllEvents godoc
// @Summary Get details of all event
// @Description Get details of all event
// @Tags event
// @Success 200 {object} event
// @Router /events [get]
func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

// updateEvent godoc
// @Summary Get details of all event
// @Description Get details of all event
// @Tags event
// @Success 200 {object} event
// @Router /events/{id} [put]
func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Wrong schema!")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}
// deleteEvent godoc
// @Summary delete details of one event
// @Description delete details of one event
// @Tags event
// @Success 200 {object} event
// @Router /events/{id} [delete]
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

// @title Events API
// @version 1.0
// @description This is a sample service for managing event
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email nguyenlephong1997@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PUT")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	// Swagger
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	//router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir("docs"))))
	log.Fatal(http.ListenAndServe(":8080", router))
}
