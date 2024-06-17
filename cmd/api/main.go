package main

import (
	"fmt"
	"context"
	"log"
	//"net/http"
    "net"

	//"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
    "google.golang.org/grpc"
    //"google.golang.org/grpc/reflection"

	"github.com/Prapul1614/RTMC/internal/user"
	"github.com/Prapul1614/RTMC/internal/rule"
    "github.com/Prapul1614/RTMC/internal/middleware"
    "github.com/Prapul1614/RTMC/proto/userpb"
    "github.com/Prapul1614/RTMC/proto/rulepb"
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