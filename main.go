package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

var client *mongo.Client

type Articles struct {
	Id                primitive.ObjectID `json:"_id" bson:"_id"`
	Title             string             `json:"title" bson:"title"`
	SubTitle          string             `json:"subtitle" bson:"subtitle"`
	Content           string             `json:"content" bson:"content"`
	CreationTimestamp time.Time          `json:"creationtimestamp" bson:"creationtimestamp"`
}

func CreateArticle(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var article Articles
	_ = json.NewDecoder(request.Body).Decode(&article)
	collection := client.Database("am").Collection("articles")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, article)
	json.NewEncoder(response).Encode(result)
}
func GetArticleById(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var article Articles
	collection := client.Database("am").Collection("articles")
	ctx, _ := context.WithTimeout(context.Background(), 40*time.Second)
	err := collection.FindOne(ctx, Articles{Id: id}).Decode(&article)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(article)
}
func ListAllArticles(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var article []Articles
	collection := client.Database("am").Collection("articles")
	ctx, _ := context.WithTimeout(context.Background(), 40*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var articles Articles
		cursor.Decode(&articles)
		article = append(article, articles)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "Message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(article)
}

func main() {
	fmt.Println("Running...")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/articles", CreateArticle).Methods("POST")
	router.HandleFunc("/articles{id}", GetArticleById).Methods("GET")
	router.HandleFunc("/articles", ListAllArticles).Methods("GET")
	http.ListenAndServe(":12345", router)
}
