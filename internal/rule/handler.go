package rule

import (
	"context"
	"encoding/json"
	"fmt"
	//"io/ioutil"
    "io"
    "errors"
    "log"
	"os"
	"net/http"
    "strings"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "google.golang.org/grpc/metadata"
    "github.com/joho/godotenv"
    "github.com/dgrijalva/jwt-go"

    "github.com/Prapul1614/RTMC/proto/rulepb"
    "github.com/Prapul1614/RTMC/internal/middleware"
)

type Handler struct {
    rulepb.UnimplementedRuleServiceServer
    service *Service
    parser *Parser
}

func NewHandler(service *Service, parser *Parser) *Handler {
    return &Handler{
        service: service,
        parser: parser,
    }
}

func (h* Handler) StreamData(stream rulepb.RuleService_StreamDataServer) error {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    var jwtKey = []byte(os.Getenv("jwtKey"))

    // Get metadata from the stream context
    md, ok := metadata.FromIncomingContext(stream.Context())
    if !ok {
        return fmt.Errorf("missing metadata")
    }

    authHead,ok := md["authorization"]
    if !ok || len(authHead) == 0 {
        return fmt.Errorf("missing authorization token")
    }

    tokenString := strings.TrimPrefix(authHead[0], "Bearer ")
    claims := &jwt.StandardClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        return fmt.Errorf("invalid token: %v", err)
    }

    owner,_ := primitive.ObjectIDFromHex(claims.Subject)
    var notifications = []string{}
    rule, err := h.service.GetRule(context.Background(), owner)
        if err != nil {
            notifications = append(notifications, "No rules added yet...")
            if err := stream.Send(&rulepb.StreamResponse{Notifications: notifications}); err != nil {
                return err
            }
        }
    for {
        msg, err := stream.Recv()
        if err == io.EOF { return nil }
        if err != nil { return err }

        textString := msg.Text // .GetText ??

        for _,v := range rule {
            ans := h.service.ImplementRule(context.Background(), textString, &v)
            if ans {
                fmt.Println("\n\nSuccess: ", v.When)
                notifications = append(notifications, v.Notify)
            } else {
                fmt.Println("\n\n Not Success: ", v.When)
            }
        }

        if err := stream.Send(&rulepb.StreamResponse{Notifications: notifications}); err != nil {
            return err
        }

    }
}

func (h *Handler) Classify(ctx context.Context, req *rulepb.ClassifyRequest) (*rulepb.ClassifyResponse, error) {
    claims, ok := ctx.Value(middleware.ClaimsKey).(*middleware.Claims)
    if !ok {
        return nil, errors.New("could not retrieve claims from context")
    }

    textString := req.Text
    
    owner,_ := primitive.ObjectIDFromHex(claims.Subject) // or claims.ID ??

    var notifications = []string{}
    rule, err := h.service.GetRule(context.Background(), owner)
    if err != nil {
        notifications = append(notifications, "No rules added yet...")
       return &rulepb.ClassifyResponse{Notifications: notifications}, nil
    }

    for _,v := range rule {
        ans := h.service.ImplementRule(context.Background(), textString, &v)
        if ans {
            fmt.Println("\n\nSuccess: ", v.When)
            notifications = append(notifications, v.Notify)
        } else {
            fmt.Println("\n\n Not Success: ", v.When)
        }
    }

    return &rulepb.ClassifyResponse{Notifications: notifications}, nil
}

func (h *Handler) CreateRule(ctx context.Context, req *rulepb.CreateRuleRequest) (*rulepb.RuleResponse, error) {    
    claims, ok := ctx.Value(middleware.ClaimsKey).(*middleware.Claims)
    if !ok {
        return nil, errors.New("could not retrieve claims from context")
    }

    ruleString := req.Rule
    owner,_ := primitive.ObjectIDFromHex(claims.Subject) 

    rule, msg, err := h.parser.ParseRule(context.Background(), ruleString, owner)
    if err != nil {
        return &rulepb.RuleResponse{Message: &msg}, nil
    }
    return &rulepb.RuleResponse{
        Rule: &rulepb.Rule{Notify: rule.Notify, When: rule.When} }, nil
}

func (h *Handler) GetRules(ctx context.Context, req *rulepb.GetRulesRequest) (*rulepb.RulesResponse, error) {    
    claims, ok := ctx.Value(middleware.ClaimsKey).(*middleware.Claims)
    if !ok {
        return nil, errors.New("could not retrieve claims from context")
    }

    owner,_ := primitive.ObjectIDFromHex(claims.Subject) 
    rules, err := h.service.GetRule(context.Background(), owner)
    if err != nil {
        msg := "Rules not found for you, Please create some and try again"
        return &rulepb.RulesResponse{Message: &msg}, nil
    }

    protoRules := make([]*rulepb.Rule, len(rules))
    for i, rule := range rules {
        protoRules[i] = &rulepb.Rule{
            Notify: rule.Notify,
            When:   rule.When,
        }
    }

    return &rulepb.RulesResponse{Rules: protoRules}, nil

}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var rule Rule
    if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if err := h.service.UpdateRule(context.Background(), id, &rule); err != nil {
        http.Error(w, "Error updating rule", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(rule)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    if err := h.service.DeleteRule(context.Background(), id); err != nil {
        http.Error(w, "Error deleting rule", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

