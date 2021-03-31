package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.Background()

type Task struct {
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("tasker").Collection("tasks")

	fmt.Println("Connected to MongoDB!!")
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	task := &Task{}

	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		fmt.Println(err)
	}

	insertResult, err := collection.InsertOne(ctx, task)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(insertResult.InsertedID)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
	}

	filer := bson.M{
		"_id": oID,
	}

	var result primitive.M

	errDecode := collection.FindOne(ctx, filer).Decode(&result)
	if errDecode != nil {
		fmt.Println("errDecode", errDecode)
	}

	json.NewEncoder(w).Encode(result)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
	}

	filer := bson.M{
		"_id": bson.M{"$eq": oID},
	}

	after := options.After
	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	newTask := &Task{}
	errDecode := json.NewDecoder(r.Body).Decode(newTask)
	if errDecode != nil {
		fmt.Println("errDecode", errDecode)
	}

	update := bson.M{
		"$set": bson.M{
			"text":      newTask.Text,
			"completed": newTask.Completed,
		},
	}

	resUpd := collection.FindOneAndUpdate(ctx, filer, update, &returnOpt)

	var result primitive.M

	errDec := resUpd.Decode(&result)
	if errDec != nil {
		fmt.Println("errDec", errDec)
	}

	json.NewEncoder(w).Encode(result)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
	}

	filer := bson.M{
		"_id": bson.M{"$eq": oID},
	}

	resDel, errDel := collection.DeleteOne(ctx, filer)
	if errDel != nil {
		fmt.Println("errDel", errDel)
	}

	json.NewEncoder(w).Encode(resDel.DeletedCount)
}

func main() {
	route := mux.NewRouter()
	s := route.PathPrefix("/api").Subrouter()

	//Routes
	s.HandleFunc("/task", createTask).Methods("POST")
	s.HandleFunc("/task/{id}", getTask).Methods("GET")
	s.HandleFunc("/task/{id}", updateTask).Methods("PUT")
	s.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", s)) // Run Server
}
