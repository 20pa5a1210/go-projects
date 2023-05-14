package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/20pa5a1210/go-projects/tree/master/mongodb-golang/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mongo.Client
}

func NewUserController(s *mongo.Client) *UserController {
	return &UserController{s}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
	}
	oid := bson.ObjectIdHex(id)
	u := models.User{}

	if err := uc.session.Database("admin").Collection("users").FindOne(r.Context(), oid); err != nil {
		w.WriteHeader(404)
		return
	}

	resp, err := json.Marshal(u)

	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s\n", &resp)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)
	// u.Id = bson.NewObjectId()
	uc.session.Database("admin").Collection("users").InsertOne(r.Context(), u)
	resp, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", resp)

}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	oid := bson.ObjectIdHex(id)
	_, err := uc.session.Database("admin").Collection("users").DeleteOne(r.Context(), oid)
	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted User", oid, "\n ")

}
