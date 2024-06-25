package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	//"net/http"
	"net"

	//"github.com/gorilla/mux"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	//"google.golang.org/grpc/reflection"

	"github.com/Prapul1614/RTMC/internal/middleware"
	"github.com/Prapul1614/RTMC/internal/rule"
	"github.com/Prapul1614/RTMC/internal/user"
	"github.com/Prapul1614/RTMC/proto/rulepb"
	"github.com/Prapul1614/RTMC/proto/userpb"
)

var client *mongo.Client

func muxserver(_ *user.Handler,rulesHandler *rule.Handler) {
    r := mux.NewRouter()
    classifyRouter := r.PathPrefix("/classify").Subrouter()
    classifyRouter.Use(middleware.JWTAuth_http)
    classifyRouter.HandleFunc("",rulesHandler.ClassifyHttp).Methods("POST")

    fmt.Println("Server started on port 8000")
    log.Fatal(http.ListenAndServe(":8000", r))
}

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
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatalf("Failed to ping MongoDB: %v",err)
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
    rulesRepo.CreateIndexOwners()

    muxserver(userHandler, rulesHandler)

    // Create a new gRPC server
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(middleware.AuthInterceptor),
    )
    

    // Register gRPC services
    userpb.RegisterUserServiceServer(grpcServer, userHandler)
    rulepb.RegisterRuleServiceServer(grpcServer, rulesHandler)

    // Enable reflection for gRPC server
    // reflection.Register(grpcServer)


    // Listen on port 3000
    lis, err := net.Listen("tcp", ":3000")
    if err != nil {
        log.Fatalf("Failed to listen on port 3000: %v", err)
    }

    log.Printf("Server is listening on port 3000...")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve gRPC server: %v", err)
    }

}