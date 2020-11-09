package main

// Imports

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"context"
	"time"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Structs and Utilities

var lock sync.Mutex // thread Safety mutex

var collection = ConnectDB()

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
    ts := time.Time(*t).Unix()
    stamp := fmt.Sprint(ts)
    return []byte(stamp), nil
}
func (t *Timestamp) UnmarshalJSON(b []byte) error {
    ts, err := strconv.Atoi(string(b))
    if err != nil {
        return err
    }
    *t = Timestamp(time.Unix(int64(ts), 0))
    return nil
}

type Article struct {
	Id		 string `json:"Id,omitempty" bson:"Id,omitempty"`
	Title    string `json:"Title,omitempty" bson:"Title,omitempty"`
	Subtitle string `json:"Subtitle,omitempty" bson:"Subtitle,omitempty"`
	Content  string `json:"Content,omitempty" bson:"Content,omitempty"`
	Creation_Timestamp *Timestamp `bson:"Creation_Timestamp,omitempty" json:"Creation_Timestamp,omitempty"`
}

// MongoDB Connection

func ConnectDB() *mongo.Collection {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017/?compressors=disabled&gssapiServiceName=mongodb")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("GolangServer").Collection("Articles")

	// MongoDB text index for searching

	/*

	indexName, err := collection.Indexes().CreateOne(
        context.Background(),
        mongo.IndexModel{
                Keys: bson.M{
					"Title":"text", 
					"Subtitle":"text",
					"Content":"text",
                },
                Options: options.Index().SetUnique(true),
        },
	)
	if err != nil {
        log.Fatal(err)
	}
	fmt.Println(indexName)
	*/

	return collection
}

type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

// getError : This is helper function to prepare error model.
func getError(err error, w http.ResponseWriter) {

	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

// Display home page
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Made by Tulasi Ram\n")
	fmt.Fprintf(w, "Portfolio: https://tulasi-ram.com\n")
	fmt.Fprintf(w, "Blog: https://blog.tulasi-ram.com\n")
	fmt.Println("Endpoint Hit: homePage")
}

// Return all articles
func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	w.Header().Set("Content-Type", "application/json")
	var Articles []Article
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		getError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var article Article

		err := cur.Decode(&article) 
		if err != nil {
			log.Fatal(err)
		}

		Articles = append(Articles, article)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(Articles)
}

// Return articles based on ID
func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnSingleArticle")
	w.Header().Set("Content-Type", "application/json")

	var article Article
	// we get params with mux.
	var params = mux.Vars(r)
	key := params["id"]
	// string to primitive.ObjectID

	filter := bson.M{"Id": key}
	err := collection.FindOne(context.TODO(), filter).Decode(&article)

	if err != nil {
		getError(err, w)
		return
	}

	json.NewEncoder(w).Encode(article)

}

// Return searched articles
func returnSearchArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnSearchArticles")
	w.Header().Set("Content-Type", "application/json")
	var Articles []Article
	//search_string := r.URL.Query().Get("q")
	search_string := "test"

	filter := bson.M{
		"$text": bson.M{
		"$search": search_string,
	   },
	}
	cur, err := collection.Find(context.TODO(), filter)
	
	if err != nil {
		getError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var article Article

		// & character returns the memory address of the following variable.
		err := cur.Decode(&article) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		Articles = append(Articles, article)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(Articles)
}

// Create new articles
func createNewArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewArticle")
	w.Header().Set("Content-Type", "application/json")

	var article Article

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&article)

	println(json.NewDecoder(r.Body).Decode(&article))
	// insert our book model.
	result, err := collection.InsertOne(context.TODO(), article)

	if err != nil {
		getError(err, w)
		return
	}
	fmt.Println("Inserted post with ID:", result.InsertedID)
	json.NewEncoder(w).Encode(result)
}

// Handle all API requests
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles).Methods("GET")
	myRouter.HandleFunc("/articles", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/articles/{id}", returnSingleArticle).Methods("GET")
	myRouter.HandleFunc("/articles/search", returnSingleArticle).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

// Main Function
func main() {
	handleRequests()
}
