package main

import (
	"fmt"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"

	"github.com/Prapul1614/RTMC/internal/user"
	"github.com/Prapul1614/RTMC/internal/rule"
)

var client *mongo.Client

func main(){
	fmt.Println("Starting API....")

	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	db := client.Database("RTMC")

	// Initialize user repository, service, and handler
    userRepo := user.NewRepository(db, "users")
    userService := user.NewService(userRepo)
    userHandler := user.NewHandler(userService)

	rulesRepo := rule.NewRepository(db, "rules")
    rulesService := rule.NewService(rulesRepo, *userRepo)
    rulesParser := rule.NewParser(rulesService)
    rulesHandler := rule.NewHandler(rulesService, rulesParser)


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

	// Add rules routes with JWT authentication middleware
    rulesRouter := r.PathPrefix("/rules").Subrouter()
    rulesRouter.Use(rule.JWTAuth)
    rulesRouter.HandleFunc("", rulesHandler.Create).Methods("POST")
    rulesRouter.HandleFunc("", rulesHandler.Get).Methods("GET")
    rulesRouter.HandleFunc("", rulesHandler.Update).Methods("PUT")
    rulesRouter.HandleFunc("", rulesHandler.Delete).Methods("DELETE")
    
    classifyRouter := r.PathPrefix("/classify").Subrouter()
    classifyRouter.Use(rule.JWTAuth)
    classifyRouter.HandleFunc("",rulesHandler.Classify).Methods("POST")

    fmt.Println("Server started on port 3000")
    log.Fatal(http.ListenAndServe(":3000", r))
}