package rule

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	//"fmt"
	"net/http"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
    service *Service
    parser *Parser
}

func NewHandler(service *Service, parser *Parser) *Handler {
    return &Handler{
        service: service,
        parser: parser,
    }
}

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

    /*var rule Rule
    if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }*/
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

