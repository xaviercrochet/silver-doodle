package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"localsearch-api/service"
)

func get(w http.ResponseWriter, r *http.Request) {
	log.Print("handling get for /place")
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	placeID := pathParams["placeID"]

	place, apiErr := service.GetPlace(placeID)

	if apiErr != nil {
		w.WriteHeader(apiErr.StatusCode)
		jsonValue, _ := json.Marshal(apiErr)
		w.Write([]byte(jsonValue))
		return
	}

	json, _ := json.Marshal(place)
	w.Write(json)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("started server on port %d", *port)
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/place/{placeID}", get).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}
