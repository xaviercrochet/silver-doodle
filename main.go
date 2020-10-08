package main

import (
	"encoding/json"
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

	//query := r.URL.Query()
	///location := query.Get("location")

	place, err := service.GetPlace(placeID)

	if err != nil {
		log.Fatalf("cannot retrieve place: %v", err)
	}

	if place == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}

	json, err := json.Marshal(place)

	if err != nil {
		log.Fatalf("cannot marshal place object to json: %v", err)
	}

	w.Write(json)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/place/{placeID}", get).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))
}
