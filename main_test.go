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

func TestGetArticles(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetArticles)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"id":1,"title":"Death Of 2","subtitle":"Death","content":"2 deaths in city","creationtimestamp":"2020-01-01 08:00:00 +0000 UTC"},{"id":2,"title":"Lockdown Lifted","subtitle":"Lockdown","content":"Lockdown has been lifted in the city","creationtimestamp":"2020-07-01 18:00:00 +0000 UTC"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
func TestGetArticleByID(t *testing.T) {

	req, err := http.NewRequest("GET", "/article", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("id", "1")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetArticleByID)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"id":1,"title":"Death Of 2","subtitle":"Death","content":"2 deaths in city","creationtimestamp":"2020-01-01 08:00:00 +0000 UTC"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
