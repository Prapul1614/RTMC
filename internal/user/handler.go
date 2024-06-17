package user

import (
    "context"

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

func (h *Handler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
    
    err := h.service.Register(ctx, req.Username, req.Password)
    if err != nil {
        return nil, err
    }
    return &userpb.RegisterResponse{Message: "User registered successfully"}, nil
}
