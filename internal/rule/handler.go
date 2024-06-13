package rule

import (
	"context"
	"encoding/json"
	//"fmt"
	"net/http"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value(claimsKey).(*Claims)
    if !ok {
        http.Error(w, "No claims found in context", http.StatusUnauthorized)
        return
    }


    var rule Rule
    if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    println(rule.Name)

    owner,_ := primitive.ObjectIDFromHex(claims.Subject)

    if _, err := h.service.FindDoc(context.Background(), &rule, owner); err == nil {
        http.Error(w, "Done", http.StatusCreated)
        return
    }

    rule.Owners = []primitive.ObjectID{}
    rule.Owners = append(rule.Owners, owner)


    if err := h.service.CreateRule(context.Background(), &rule, owner); err != nil {
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

    /*id, err := primitive.ObjectIDFromHex(r.URL.Query().Get("id"))
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }*/

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
