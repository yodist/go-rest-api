package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yodist/go-rest-api/config"
	"github.com/yodist/go-rest-api/models"
)

var conf = config.Config{}
var client *mongo.Client

func allMoviesEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var movies []models.Movie
	collection := client.Database(conf.Database).Collection("movie")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var movie models.Movie
		cursor.Decode(&movie)
		movies = append(movies, movie)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(movies)
}

func findMovieEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var movie models.Movie
	collection := client.Database(conf.Database).Collection("movie")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, models.Movie{ID: id}).Decode(&movie)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(movie)
}

func createMovieEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var movie models.Movie
	_ = json.NewDecoder(request.Body).Decode(&movie)
	collection := client.Database(conf.Database).Collection("movie")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, movie)
	json.NewEncoder(response).Encode(result)
}

func updateMovieEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var movie models.Movie
	_ = json.NewDecoder(request.Body).Decode(&movie)
	collection := client.Database(conf.Database).Collection("movie")
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", movie},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(response).Encode(result)
}

func deleteMovieEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func main() {

	// run pre-configuration needed
	fmt.Println("served on port :3000")
	conf.Read()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(conf.Server))

	// set endpoint
	r := mux.NewRouter()
	r.HandleFunc("/movies", allMoviesEndPoint).Methods("GET")
	r.HandleFunc("/movies", createMovieEndPoint).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovieEndPoint).Methods("PUT")
	r.HandleFunc("/movies", deleteMovieEndPoint).Methods("DELETE")
	r.HandleFunc("/movies/{id}", findMovieEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
