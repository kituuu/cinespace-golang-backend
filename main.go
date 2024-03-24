package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Username      string   `json:"username"`
	Name          string   `json:"name"`
	Subscriptions []string `json:"subs"`
	SubsAdd       string   `json:"subsAdd"`
	History       []string `json:"history"`
	HistoryAdd    string   `json:"hisAdd"`
	VideoUploaded []string `json:"vidUpload"`
	VidUplAdd     string   `json:"vidUplAdd"`
	Avatar        string   `json:"avatar"`
	TotalViews    int      `json:"totalviews"`
	IncreaseViews bool     `json:"incViews"`
}

var Users []User

var usersCollection *mongo.Collection

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// code for fetch all users
	fmt.Println("Endpoint Hit: getUsers")
	cursor, err := usersCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	fmt.Println("displaying all results in a collection")
	for _, result := range results {
		fmt.Println(result)
	}
	json.NewEncoder(w).Encode(results)
}

func updateUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := (r.URL.Query())
	key := vars["username"][0]
	fmt.Println("Endpoint Hit: updateUsers")
	cursor, err := usersCollection.Find(context.TODO(), bson.M{"username": key})
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	if results == nil {
		fmt.Fprintln(w, "No such User exist. Create a User")
	} else {
		reqBody, _ := io.ReadAll(r.Body)
		var newUser, result User
		var present bool
		json.Unmarshal(reqBody, &newUser)
		by, _ := json.Marshal(results[0])
		json.Unmarshal(by, &result)
		fmt.Println("resul", result.History)
		fmt.Println(results[0])
		fmt.Println(newUser)
		fmt.Println(bson.A{newUser}[0])
		if newUser.Avatar == "" {
			newUser.Avatar = result.Avatar
		}
		if newUser.Name == "" {
			newUser.Name = result.Name
		}
		if newUser.HistoryAdd == "" {
			newUser.History = result.History
		} else {
			present = false
			for _, res := range result.History {
				fmt.Println(res)
				if res == newUser.HistoryAdd {
					present = true
				}
			}
			if !present {
				newUser.History = append(result.History, newUser.HistoryAdd)
			} else {
				newUser.History = result.History
			}

		}
		if newUser.SubsAdd == "" {
			newUser.Subscriptions = result.Subscriptions
		} else {
			present = false
			for _, res := range result.Subscriptions {
				fmt.Println(res)
				if res == newUser.SubsAdd {
					present = true
				}
			}
			if !present {
				newUser.Subscriptions = append(result.Subscriptions, newUser.SubsAdd)
			} else {
				newUser.Subscriptions = result.Subscriptions
			}
		}
		if newUser.VidUplAdd == "" {
			newUser.VideoUploaded = result.VideoUploaded
		} else {
			present = false
			for _, res := range result.VideoUploaded {
				fmt.Println(res)
				if res == newUser.VidUplAdd {
					present = true
				}
			}
			if !present {
				newUser.VideoUploaded = append(result.VideoUploaded, newUser.VidUplAdd)
			} else {
				newUser.VideoUploaded = result.VideoUploaded
			}
		}
		if newUser.IncreaseViews {
			newUser.TotalViews = 1 + result.TotalViews
		} else {
			newUser.TotalViews = result.TotalViews
		}
		fmt.Println(newUser.TotalViews)
		_, err := usersCollection.UpdateOne(context.TODO(), bson.M{"username": key}, bson.D{{Key: "$set", Value: bson.D{{Key: "subs", Value: newUser.Subscriptions}, {Key: "avatar", Value: newUser.Avatar}, {Key: "history", Value: newUser.History}, {Key: "vidUpload", Value: newUser.VideoUploaded}, {Key: "name", Value: newUser.Name}, {Key: "totalviews", Value: newUser.TotalViews}}}})
		if err != nil {
			panic(err)
		}
		getUsers(w, r)
	}

}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	reqBody, _ := io.ReadAll(r.Body)
	var newUser User
	json.Unmarshal(reqBody, &newUser)
	//code for creating
	user := bson.D{{Key: "username", Value: newUser.Username}, {Key: "name", Value: newUser.Name}, {Key: "subs", Value: newUser.Subscriptions}, {Key: "history", Value: newUser.History}, {Key: "vidUpload", Value: newUser.VideoUploaded}, {Key: "avatar", Value: newUser.Avatar}, {Key: "totalviews", Value: newUser.TotalViews}}

	result, err := usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	getUsers(w, r)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := (r.URL.Query())
	key := vars["username"][0]
	fmt.Println("Endpoint Hit: getUser", key)
	//code for filter user
	cursor, err := usersCollection.Find(context.TODO(), bson.M{"username": key})
	if err != nil {
		panic(err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	fmt.Println("displaying all results in a collection")
	for _, result := range results {
		fmt.Println(result)
	}
	if results != nil {
		json.NewEncoder(w).Encode(results)
	} else {
		fmt.Fprintln(w, "No such data")
	}

}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/users", getUsers)
	router.HandleFunc("/user", updateUsers).Methods("PUT")
	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/user", getUser)

	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	fmt.Println("Listening to http://localhost:10000/")
	godotenv.Load()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	usersCollection = client.Database("Cinespace").Collection("users")
	fmt.Println(client)
	handleRequests()
}
