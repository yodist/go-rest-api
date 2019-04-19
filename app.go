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
		response.WriteHeader(http.StatusBadRequest)
		res := models.Response{Status: http.StatusBadRequest, StatusCode: 1, Message: err.Error()}
		json.NewEncoder(response).Encode(res)
		log.Println(err)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var movie models.Movie
		cursor.Decode(&movie)
		movies = append(movies, movie)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		res := models.Response{Status: http.StatusBadRequest, StatusCode: 1, Message: err.Error()}
		json.NewEncoder(response).Encode(res)
		log.Println(err)
		return
	}
	res := models.Response{Data: movies, Status: http.StatusOK, StatusCode: 0, Message: "Success"}
	json.NewEncoder(response).Encode(res)
}

func findMovieEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var movie models.Movie
	collection := client.Database(conf.Database).Collection("movie")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&movie)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		res := models.Response{Status: http.StatusBadRequest, StatusCode: 1, Message: err.Error()}
		json.NewEncoder(response).Encode(res)
		log.Println(err)
		return
	}

	res := models.Response{Data: movie, Status: http.StatusOK, StatusCode: 0, Message: "Success"}
	json.NewEncoder(response).Encode(res)
}

func createMovieEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var movie models.Movie
	err := json.NewDecoder(request.Body).Decode(&movie)
	collection := client.Database(conf.Database).Collection("movie")
	if movie.CreatedBy == "" {
		movie.CreatedBy = "system"
	}
	timeNow := time.Now()
	movie.CreatedDate = &timeNow
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = collection.InsertOne(ctx, movie)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		res := models.Response{Status: http.StatusBadRequest, StatusCode: 1, Message: err.Error()}
		json.NewEncoder(response).Encode(res)
		log.Println(err)
		return
	}
	res := models.Response{Data: movie, Status: http.StatusOK, StatusCode: 0, Message: "Success"}
	json.NewEncoder(response).Encode(res)
}

func createMultipleMovieEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var movies []models.Movie
	err := json.NewDecoder(request.Body).Decode(&movies)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		res := models.Response{Status: http.StatusBadRequest, StatusCode: 1, Message: err.Error()}
		json.NewEncoder(response).Encode(res)
		log.Println(err)
		return
	}
	var moviesInterface []interface{}
	for _, movie := range movies {
		if movie.CreatedBy == "" {
			movie.CreatedBy = "system"
		}
		timeNow := time.Now()
		movie.CreatedDate = &timeNow
		moviesInterface = append(moviesInterface, movie)
	}
	collection := client.Database(conf.Database).Collection("movie")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertMany(ctx, moviesInterface)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		res := models.Response{Status: http.StatusBadRequest, StatusCode: 1, Message: err.Error()}
		json.NewEncoder(response).Encode(res)
		log.Println(err)
		return
	}
	res := models.Response{Data: movies, Status: http.StatusOK, StatusCode: 0, Message: "Success"}
	json.NewEncoder(response).Encode(res)
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
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		res := models.Response{Status: http.StatusBadRequest, StatusCode: 1, Message: err.Error()}
		json.NewEncoder(response).Encode(res)
		log.Println(err)
		return
	}
	res := models.Response{Data: movie, Status: http.StatusOK, StatusCode: 0, Message: "Success"}
	json.NewEncoder(response).Encode(res)
}

func updateMultipleMovieEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var movies []models.Movie
	_ = json.NewDecoder(request.Body).Decode(&movies)

	// collection := client.Database(conf.Database).Collection("movie")
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// var result

	// for _, movie := range movies {
	// 	filter := bson.D{{"_id", movie.ID}}
	// 	update := bson.D{
	// 		{"$set", movie},
	// 	}
	// 	result, err := collection.UpdateOne(ctx, filter, update)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }

	json.NewEncoder(response).Encode("Not Implemented")
}

func deleteMovieEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database(conf.Database).Collection("movie")
	filter := bson.D{{"_id", id}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		res := models.Response{Status: http.StatusBadRequest, StatusCode: 1, Message: err.Error()}
		json.NewEncoder(response).Encode(res)
		log.Println(err)
		return
	}
	res := models.Response{Data: result, Status: http.StatusOK, StatusCode: 0, Message: "Success"}
	json.NewEncoder(response).Encode(res)
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
	r.HandleFunc("/movies/{id}", findMovieEndpoint).Methods("GET")
	//r.HandleFunc("/movies", createMovieEndPoint).Methods("POST")
	r.HandleFunc("/movies", createMultipleMovieEndPoint).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovieEndPoint).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovieEndPoint).Methods("DELETE")

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
