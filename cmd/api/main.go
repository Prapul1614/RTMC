package main

import (
	"fmt"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Prapul1614/RTMC/internal/user"
)

var client *mongo.Client

func main(){
	fmt.Println("Starting API....")

	var err error
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/RTMC")
    client, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")

	// Initialize user repository, service, and handler
    userRepo := user.NewRepository(client.Database("RTMC"), "users")
    userService := user.NewService(userRepo)
    userHandler := user.NewHandler(userService)

	// Set up router
    r := mux.NewRouter()
	// Add a simple health check endpoint
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Hello, World!"))
    }).Methods("GET")
    

	// Add user routes
    r.HandleFunc("/login", userHandler.Login).Methods("POST")
    r.HandleFunc("/register", userHandler.Register).Methods("POST")

    fmt.Println("Server started on port 3000")
    log.Fatal(http.ListenAndServe(":3000", r))
}