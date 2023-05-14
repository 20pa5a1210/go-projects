package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"first_name" bson:"first_name,omitempty"`
	LastName  string             `json:"last_name" bson:"last_name,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
}

var client *mongo.Client

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}
	collection := client.Database("test").Collection("users")
	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
		return
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()
	id, err := primitive.ObjectIDFromHex(params.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}
	collection := client.Database("test").Collection("users")
	var user User
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}
	json.NewEncoder(w).Encode(user)
}

func main() {
	// Set up MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://admin:admin@cluster0.jhj7ej8.mongodb.net/test"))
	if err != nil {
		log.Fatal(err)
	}

	// Define routes
	router := http.NewServeMux()
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createUserHandler(w, r)
		case http.MethodGet:
			getUserHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		}
	})

	// Start server
	log.Fatal(http.ListenAndServe(":8080", router))
}
