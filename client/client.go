package main

import (
	"context"
	"log"
	"time"
	"math/rand"

	"github.com/Prapul1614/RTMC/proto/rulepb"
	"github.com/Prapul1614/RTMC/proto/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)
var tokens = []string{
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg3MDkwMDgsInN1YiI6IjY2NzAxOGFhMDc4YzIwODc0ODQwNWRlYSJ9.uoOVjD4WYLV--PYQlVdFdoxTTvNwhK4kV16S3_EksTE",
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg3MDkwNDYsInN1YiI6IjY2NzAxOGM1MDc4YzIwODc0ODQwNWRlYiJ9.tXdOW38EmF0bGMpcog5wO5JxHx6Vd7B1lBxd64bbLe8",
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg3MDkwODMsInN1YiI6IjY2NzAxOGQ0MDc4YzIwODc0ODQwNWRlYyJ9.1rM4Y78OMKm2SOkpdQN_awdSE3d0mKkHU_qIPcs0JpA",
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg3MDkxMTksInN1YiI6IjY2NzAxOGRmMDc4YzIwODc0ODQwNWRlZCJ9.jpo7K61rEdLpRpOFssSoJZWs7uFEzdUxJdCD-h0WatQ",
}

var texts = []string{
	"transaction transaction location",
	"urgent urgent money transfer",
	"stock price now decrease",
	"three failes login attempts",
	"delay traffic delay delay",
	" you are mentions in and good for mentions",
	"energy high very very",
	"your request for product recall successfull",
	"Yoo its going to be an earthquake",
	"hospital patients waiting patients no rooms for patients",
}

func TestRegister(client userpb.UserServiceClient) {
	// Test Register method
    registerRequest := &userpb.RegisterRequest{
        Username: "stream_user4",
        Password: "stream_password4",
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
        Username: "stream_user4",
        Password: "stream_password4",
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
	token := tokens[2]
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	createRequest := &rulepb.CreateRuleRequest{
		Rule: "NOTIFY High patient volume warning WHEN AND (Contains \"hospital\" , MAX (Count \"patients\" , Count \"waiting\") >= 3)",
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

func TestStream(client rulepb.RuleServiceClient) {
	token := tokens[0]
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	stream, err := client.StreamData(ctx)
    if err != nil {
        log.Fatalf("could not create stream: %v", err)
    }
	for i := 0;i < 10;i++ {
		err := stream.Send(&rulepb.StreamRequest{Text: texts[rand.Intn(10)]})
        if err != nil {
            log.Fatalf("could not send text: %v", err)
        }

        // Receive response from server
        response, err := stream.Recv()
        if err != nil {
            log.Fatalf("could not receive response: %v", err)
        }
        log.Printf("Received:")
		for _,v := range response.Notifications {
			log.Print(v)
		}
    }

    stream.CloseSend()
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

	//TestClassify(client)

	TestStream(client)

	//fmt.Print(client)

}
