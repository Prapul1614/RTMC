package test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	"github.com/Prapul1614/RTMC/internal/middleware"
	"github.com/Prapul1614/RTMC/internal/rule"
	"github.com/Prapul1614/RTMC/internal/user"
	"github.com/Prapul1614/RTMC/proto/rulepb"
	"github.com/Prapul1614/RTMC/proto/userpb"
)

func Generate(t *testing.T) string {
	rand.Seed(time.Now().UnixNano())

	var letterRunes []rune
	for i := 'a'; i <= 'z'; i++ {
		letterRunes = append(letterRunes, i)
	}
	for i := 'A'; i <= 'Z'; i++ {
		letterRunes = append(letterRunes, i)
	}

	b := make([]rune, 6)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
var username, password string
var jwtKey []byte
func Set(t *testing.T) {
	username = Generate(t)
	password = Generate(t)
}

func newServer(t *testing.T) (*grpc.ClientConn, *user.Repository, *rule.Repository) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/TestRTMC")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

	db := client.Database("TestRTMC")

	// Initialize user repository, service, and handler
    userRepo := user.NewRepository(db, "users")
    userService := user.NewService(userRepo)
    userHandler := user.NewHandler(userService)

	rulesRepo := rule.NewRepository(db, "rules")
    rulesService := rule.NewService(rulesRepo, *userRepo)
    rulesParser := rule.NewParser(rulesService)
    rulesHandler := rule.NewHandler(rulesService, rulesParser)

	grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(middleware.AuthInterceptor),
    )
	userpb.RegisterUserServiceServer(grpcServer, userHandler)
	rulepb.RegisterRuleServiceServer(grpcServer, rulesHandler)

	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	t.Cleanup(func() {
		grpcServer.Stop()
	})

	err = godotenv.Load()
    if err != nil {
        t.Fatalf("Error loading .env file: %v", err)
    }
	jwtKey = []byte(os.Getenv("jwtKey"))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("grpcServer.Serve %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatalf("grpc.DialContext %v", err)
	}
	return conn, userRepo, rulesRepo
}

func TestRegister(t *testing.T) {
	conn, userRepo, _ := newServer(t)

	client := userpb.NewUserServiceClient(conn)

	// username := "username1"
	// password := "password1"
	Set(t)
	registerRequest := &userpb.RegisterRequest{
        Username: username,
        Password: password,
    }
	_, err := client.Register(context.Background(), registerRequest)
	if err != nil {
		t.Fatalf("client.Register \n %v message ", err)
	}

	var user user.User
	filter := bson.M{"username": username}
    err = userRepo.Collection().FindOne(context.Background(), filter).Decode(&user)
    if err != nil {
        t.Fatalf("Can't Find Registered user: %v", err)
    }
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        t.Fatalf("Extracted password is different")
    }
}

func TestLogin(t *testing.T) {
	conn, userRepo, _ := newServer(t)
	client := userpb.NewUserServiceClient(conn)

	LoginRequest := &userpb.LoginRequest{
        Username: username,
        Password: password + ".",
    }
    _, err := client.Login(context.Background(), LoginRequest)
	if err == nil {
		t.Fatalf("client.Login \n %v message, should have given \"invalid credential\" error ", err)
	}

	LoginRequest = &userpb.LoginRequest{
        Username: username,
        Password: password,
    }
    res, err := client.Login(context.Background(), LoginRequest)
	if err != nil {
		t.Fatalf("client.Login \n %v message ", err)
	}
	tokenString := *res.Token

	claims := &jwt.StandardClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        t.Fatalf("invalid token: %v", err)
    }
    owner,err := primitive.ObjectIDFromHex(claims.Subject)
	if err != nil || !token.Valid {
        t.Fatalf("unable to convert to id from string: %v", err)
    }

	var user user.User
	filter := bson.M{"_id": owner}
    err = userRepo.Collection().FindOne(context.Background(), filter).Decode(&user)
    if err != nil {
        t.Fatalf("Can't Find user: %v", err)
    }
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        t.Fatalf("Extracted password is different")
    }	
}

func TestCreateRule(t *testing.T) {
	conn , _ , rulesrepo := newServer(t)

	LoginRequest := &userpb.LoginRequest{
        Username: "username0",
        Password: "password0",
    }
	client := userpb.NewUserServiceClient(conn)
    res, err := client.Login(context.Background(), LoginRequest)
	if err != nil {
		t.Fatalf("client.Login \n %v message ", err)
	}
	token := *res.Token
	
	client1 := rulepb.NewRuleServiceClient(conn)
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, md)

	//Notify := "High patient volume warning"
	//When := "AND (Contains \"hospital\" , MAX (Count \"patients\" , Count \"waiting\") >= 3)"

	Notify := "Too many a's"
	When := "Count \"a\" > 5"

	createRequest := &rulepb.CreateRuleRequest{
		Rule: fmt.Sprintf("NOTIFY %v WHEN %v", Notify, When),
	}
	rule_res , err := client1.CreateRule(ctx, createRequest)
	if err != nil {
        t.Fatalf("Error while calling CreateRule RPC: %v", err)
    }
	if rule_res.Rule.Notify != Notify || rule_res.Rule.When != When {
		t.Fatalf("Created Rule classification message or condition is not matching.")
	}

	filter := bson.M{ "notify": Notify, "when" : When,}
    var temp rule.Rule
    err = rulesrepo.Collection().FindOne(context.TODO(),filter).Decode(&temp)
	if err != nil {
		t.Fatalf("Rule Not Found: %v", err)
	}

	claims := &jwt.StandardClaims{}
    jwttoken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !jwttoken.Valid {
        t.Fatalf("invalid token: %v", err)
    }
    owner,err := primitive.ObjectIDFromHex(claims.Subject)
	if err != nil {
        t.Fatalf("unable to convert to id from string: %v", err)
    }
	var found bool 
	for _,v := range(temp.Owners) {
		if v == owner {found = true; break}
	}
	if !found {
		t.Fatalf("Owner not found in rule owners list!")
	}
}

func TestClassify(t *testing.T) {
	conn, _ , _ := newServer(t)
	LoginRequest := &userpb.LoginRequest{
        Username: "username0",
        Password: "password0",
    }
	client := userpb.NewUserServiceClient(conn)
    res, err := client.Login(context.Background(), LoginRequest)
	if err != nil {
		t.Fatalf("client.Login \n %v message ", err)
	}
	token := *res.Token
	
	client1 := rulepb.NewRuleServiceClient(conn)
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, md)

	var strings = []string{
		"waiting.. hospital is very busy so patients are waiting, waiting..",
		"kjdfn aeuoifncm enfiom inmcm",
		"waiting patients in hospital",
		"aaaaaaaaaaaaa",
	}

	// 2 rules 
	// 1. NOTIFY High patient volume warning WHEN AND (Contains \"hospital\" , MAX (Count \"patients\" , Count \"waiting\") >= 3)
	// 2. NOTIFY Too many a's WHEN Count "a" > 5
	for i,v := range(strings) {
		classifyRequset := &rulepb.ClassifyRequest {Text: v}
		res, err := client1.Classify(ctx, classifyRequset)
		if err != nil {
			t.Fatalf("client.Classify: %v", err)
		}
		if i == 0 {
			if len(res.Notifications) != 2 {
				t.Fatalf("Should have got 2 notifications")
			}
		}else if i == 1 || i == 2 {
			if len(res.Notifications) != 0 {
				t.Fatalf("Should have not got any  notifications")
			}
		}else {
			if len(res.Notifications) !=  1 || res.Notifications[0] != "Too many a's"{
				t.Fatalf("Should have got only \"Too many a's\" notification")
			}
		}
	}
}