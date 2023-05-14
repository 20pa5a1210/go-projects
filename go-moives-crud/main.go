package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Moive struct {
	Id       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}
type Director struct {
	FirstName string `json:"firstname"`
	LasName   string `json:"lastname"`
}

var moives []Moive

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "applicatio/json")
	json.NewEncoder(w).Encode(moives)
}

func deleteMoive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "applicatio/json")
	params := mux.Vars(r)
	for index, item := range moives {
		if item.Id == params["id"] {
			moives = append(moives[:index], moives[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(moives)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "applicatio/json")
	params := mux.Vars(r)
	for _, item := range moives {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

}

func createMoive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "applicatio/json")
	var moive Moive
	_ = json.NewDecoder(r.Body).Decode(&moive)
	moive.Id = strconv.Itoa(rand.Intn(1000000))

	moives = append(moives, moive)
	json.NewEncoder(w).Encode(moive)
}
func updateMoive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "applicatio/json")
	params := mux.Vars(r)
	for index, item := range moives {
		if item.Id == params["id"] {
			moives = append(moives[:index], moives[index+1:]...)
			var moive Moive
			_ = json.NewDecoder(r.Body).Decode(&moive)
			moive.Id = params["id"]
			moives = append(moives, moive)
			json.NewEncoder(w).Encode(moive)
			return
		}
	}
}

func main() {

	moives = append(moives,
		Moive{
			Id:    "1",
			Isbn:  "75792",
			Title: "Batman Begins",
			Director: &Director{
				FirstName: "Christopher",
				LasName:   "nolan",
			}},
		Moive{
			Id:    "2",
			Isbn:  "875792",
			Title: "Batman Dark-Knight",
			Director: &Director{
				FirstName: "Christopher",
				LasName:   "nolan",
			}},
	)

	router := mux.NewRouter()
	router.HandleFunc("/moives", getMovies).Methods("GET")
	router.HandleFunc("/moives/{id}", getMovie).Methods("GET")
	router.HandleFunc("/moives", createMoive).Methods("POST")
	router.HandleFunc("/moives/{id}", updateMoive).Methods("PUT")
	router.HandleFunc("/moives/{id}", deleteMoive).Methods("DELETE")

	fmt.Println("Server up and running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
