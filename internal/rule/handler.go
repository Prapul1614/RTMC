package rule

import (
	"context"
	"encoding/json"
	"fmt"
	//"io/ioutil"
    "errors"
	//"fmt"
	"net/http"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
/*
func (h *Handler) Classify(w http.ResponseWriter, r *http.Request) {

    claims, ok := r.Context().Value(claimsKey).(*Claims)
    if !ok {
        http.Error(w, "No claims found in context", http.StatusUnauthorized)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    textString := string(body)
    fmt.Println("\n\n\n\n",textString)

    owner,_ := primitive.ObjectIDFromHex(claims.Subject)

    rule, err := h.service.GetRule(context.Background(), owner)
    if err != nil {
        http.Error(w, "Rules not found", http.StatusNotFound)
        return
    }

    for _,v := range rule {
        ans := h.service.ImplementRule(context.Background(), textString, &v)
        if ans {
            fmt.Println("\n\nSuccess: ", v.When)
        } else {
            fmt.Println("\n\n Not Success: ", v.When)
        }
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value(claimsKey).(*Claims)
    if !ok {
        http.Error(w, "No claims found in context", http.StatusUnauthorized)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    ruleString := string(body)
    fmt.Println("\n\n\n\n",ruleString)

    owner,_ := primitive.ObjectIDFromHex(claims.Subject)

    rule, err := h.parser.ParseRule(context.Background(), ruleString, owner)
    if err != nil {
        println(err.Error())
        http.Error(w, "Error creating rule", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(rule)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value(claimsKey).(*Claims)
    if !ok {
        http.Error(w, "No claims found in context", http.StatusUnauthorized)
        return
    }

    owner,_ := primitive.ObjectIDFromHex(claims.Subject)

    rule, err := h.service.GetRule(context.Background(), owner)
    if err != nil {
        http.Error(w, "Rule not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(rule)
}*/

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

