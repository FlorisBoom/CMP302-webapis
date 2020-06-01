package main

import (
	"log"
	"encoding/json"
	"golang.org/x/net/context"
	"net/http"
	"github.com/gorilla/mux"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	// "google.golang.org/api/option"
	"os"
)

// Car Struct
type Car struct {
	Brand	string `firestore:"Brand,omitempty"`
	Model	string `firestore:"Model,omitempty"`
	Year	int64  `firestore:"Year,omitempty"`
	Color	string `firestore:"Color,omitempty"`
}

// var opt = option.WithCredentialsFile("./firebase-key.json")
var app, err = firebase.NewApp(context.Background(), nil)

func main() {
	// Init firebase
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Init Router
	r := mux.NewRouter()

	// Route Handlers for endpoints
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/bin/index.html")
	})
	r.HandleFunc("/cars", getCars).Methods("GET")
	r.HandleFunc("/car/{id}", getCar).Methods("GET")
	r.HandleFunc("/cars", createCar).Methods("POST")
	r.HandleFunc("/cars/{id}", updateCar).Methods("PUT")	
	r.HandleFunc("/cars/{id}", deleteCar).Methods("DELETE")

	// Auth Request for token
	r.HandleFunc("/authorize", getAuthorization).Methods("GET")


	port := getPort()
	log.Fatal(http.ListenAndServe(port, r))
}

func getCars(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if !verifyToken(token) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401"))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")

	cars := make(map[string] Car)

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer client.Close()

	iter := client.Collection("cars").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		car := Car {
			Brand: doc.Data()["Brand"].(string),
			Model: doc.Data()["Model"].(string),
			Year:  doc.Data()["Year"].(int64),
			Color: doc.Data()["Color"].(string),
		}
		cars[doc.Ref.ID] = car
	}
	json.NewEncoder(w).Encode(cars)
}

func getCar(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if !verifyToken(token) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer client.Close()

	params := mux.Vars(r)

	dsnap, err := client.Collection("cars").Doc(params["id"]).Get(context.Background())
	if err != nil {
		log.Fatalln(err)
		return
	}

	car := dsnap.Data()
	json.NewEncoder(w).Encode(car)
}

func createCar(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if !verifyToken(token) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer client.Close()
	 
	var car Car
	_ = json.NewDecoder(r.Body).Decode(&car)

	_, _, err = client.Collection("cars").Add(context.Background(), car)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}

	json.NewEncoder(w).Encode(car)
}

func updateCar(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if !verifyToken(token) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer client.Close()

	params := mux.Vars(r)

	var car Car 
	_ = json.NewDecoder(r.Body).Decode(&car)

	_, err = client.Collection("cars").Doc(params["id"]).Set(context.Background(), car)
	if err != nil {
			log.Printf("An error has occurred: %s", err)
	}

	json.NewEncoder(w).Encode(car)
}

func deleteCar(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if !verifyToken(token) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401"))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)                                  
		return
	}
	defer client.Close()
	 
	params := mux.Vars(r)

	var car Car 
	_ = json.NewDecoder(r.Body).Decode(&car)

	_, err = client.Collection("cars").Doc(params["id"]).Delete(context.Background())
	if err != nil {
			log.Printf("An error has occurred: %s", err)
	}
}

func getAuthorization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	client, err := app.Auth(context.Background())
	if err != nil {
			log.Fatalf("error getting Auth client: %v\n", err)
	}
	 

	token, err := client.CustomToken(context.Background(), "some-uid")
	if err != nil {
			log.Fatalf("error minting custom token: %v\n", err)
	}

	uploadToken(token)
	log.Printf("Got custom token: %v\n", token)
	json.NewEncoder(w).Encode(token)

}

func uploadToken(token string) {
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)                                  
		return
	}

	_, _, err = client.Collection("tokens").Add(context.Background(), map[string]interface{}{
		"token": token,
	})
	if err != nil {
		log.Fatalf("An error has occurred: %s", err)
	}
}

func verifyToken(token string) bool {
	println(token)
	if token == "" {
		return false
	} 

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)                                  
		return false
	}

	iter := client.Collection("tokens").Where("token", "==", token).Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("An error has occurred: %s", err)
			return false
		}
		println(doc.Data())
	}
	log.Printf("Verified token: %v\n", token)
	return true
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
	  return ":" + p
	}
	return ":8000"
}