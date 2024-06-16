package user

import (
    "context"
    //"encoding/json"
    //"net/http"

    "github.com/Prapul1614/RTMC/proto/userpb"
)

type Handler struct {
    userpb.UnimplementedUserServiceServer
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
    
    token, err := h.service.Authenticate(ctx, req.Username, req.Password)
    if err != nil {
        return nil, err
    }
    return &userpb.LoginResponse{Token: &token}, nil
}
/*func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    token, err := h.service.Authenticate(context.Background(), req.Username, req.Password)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    err := h.service.Register(context.Background(), req.Username, req.Password)
    if err != nil {
        http.Error(w, "Error registering user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}*/

func (h *Handler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
    
    err := h.service.Register(ctx, req.Username, req.Password)
    if err != nil {
        return nil, err
    }
    return &userpb.RegisterResponse{Message: "User registered successfully"}, nil
}
