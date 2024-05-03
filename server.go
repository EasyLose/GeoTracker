package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type User struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type POI struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var users = make(map[string]User)
var pois = []POI{
	{"Central Park", 40.785091, -73.968285},
	{"Empire State Building", 40.748817, -73.985428},
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func updateUserLocation(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users[user.ID] = user
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Location updated successfully")
}

func getNearbyPOIs(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nearbyPOIs := make([]POI, 0)
	for _, poi := range pois {
		if abs(user.Latitude-poi.Latitude) < 0.1 && abs(user.Longitude-poi.Longitude) < 0.1 {
			nearbyPOIs = append(nearbyPOIs, poi)
		}
	}

	response, _ := json.Marshal(nearbyPOIs)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}

func main() {
	loadEnv()

	router := mux.NewRouter()

	apiRoute := router.PathPrefix("/api").Subrouter()
	apiRoute.HandleFunc("/updateLocation", updateUserLocation).Methods("POST")
	apiRoute.HandleFunc("/nearbyPOIs", getNearbyPOIs).Methods("POST")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:" + os.Getenv("PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}