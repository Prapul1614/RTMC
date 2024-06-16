package main

import (
	"context"
	"log"
	"time"

	"github.com/Prapul1614/RTMC/proto/rulepb"
	"github.com/Prapul1614/RTMC/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestRegister(client userpb.UserServiceClient) {
	// Test Register method
	println("Enterining")
    registerRequest := &userpb.RegisterRequest{
        Username: "example_user1",
        Password: "example_password1",
    }
    registerResponse, err := client.Register(context.Background(), registerRequest)
    if err != nil {
        log.Fatalf("Error while calling Register RPC: %v", err)
    }
    log.Printf("Register Response: %v", registerResponse)
}

func TestLogin(client userpb.UserServiceClient) {
	// Test Login method
    loginRequest := &userpb.LoginRequest{
        Username: "example_user1",
        Password: "example_password1",
    }
    loginResponse, err := client.Login(context.Background(), loginRequest)
    if err != nil {
        log.Fatalf("Error while calling Login RPC: %v", err)
    }
    log.Printf("Login Response: %v", loginResponse)
}

func TestCreate(client rulepb.RuleServiceClient) {
	// Create a new context with a timeout
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

	// Add authentication token to the context
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg2MDUxNDgsInN1YiI6IjY2NmU4MTU5OTVkYmIxYzQ1MTMyMzUwMyJ9.-QQbLkRGQcQFkc5_YpDHDamj0givVU2cqcJJcuGU9Xk" // Replace with the actual token
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	createRequest := &rulepb.CreateRuleRequest{
		Rule: "NOTIFY hii WHEN AND (MAX (Count \"pop\" , Length) < 2 , Contains \"shit\")",
	}
	createResponse, err := client.CreateRule(ctx, createRequest)
	if err != nil {
        log.Fatalf("Error while calling CreateRule RPC: %v", err)
    }
    log.Printf("Login Response: %v", createResponse)
}

func TestGet(client rulepb.RuleServiceClient) {
	// Create a new context with a timeout
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

	// Add authentication token to the context
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg2MDUxNDgsInN1YiI6IjY2NmU4MTU5OTVkYmIxYzQ1MTMyMzUwMyJ9.-QQbLkRGQcQFkc5_YpDHDamj0givVU2cqcJJcuGU9Xk" // Replace with the actual token
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	getRequest := &rulepb.GetRulesRequest{}
	getResponse, err := client.GetRules(ctx, getRequest)
	
	if err != nil {
        log.Fatalf("Error while calling CreateRule RPC: %v", err)
    }
    log.Printf("Login Response: %v", getResponse)
}

func TestClassify(client rulepb.RuleServiceClient) {
	// Create a new context with a timeout
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

	// Add authentication token to the context
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg2MDUxNDgsInN1YiI6IjY2NmU4MTU5OTVkYmIxYzQ1MTMyMzUwMyJ9.-QQbLkRGQcQFkc5_YpDHDamj0givVU2cqcJJcuGU9Xk" // Replace with the actual token
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	classifyRequest := &rulepb.ClassifyRequest{
		Text: "pop",
	}
	classifyResponse, err := client.Classify(ctx, classifyRequest)
	if err != nil {
        log.Fatalf("Error while calling CreateRule RPC: %v", err)
    }
    log.Printf("Login Response: %v", classifyResponse)
}

func main() {
    conn, err := grpc.Dial("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Did not connect: %v", err)
    }
    defer conn.Close()

    //client := userpb.NewUserServiceClient(conn)
	client := rulepb.NewRuleServiceClient(conn)

    //TestRegister(client)

    //TestLogin(client)

	//TestCreate(client)

	//TestGet(client)

	TestClassify(client)

	//fmt.Print(client)

}
