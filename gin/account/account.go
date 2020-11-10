package account

import (
	context "context"

	grpc "google.golang.org/grpc"
)

type AccountService struct {
}

func (a *AccountService) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	return nil, nil
}
