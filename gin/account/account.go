package account

import (
	context "context"
)

type AccountService struct {
}

func (a *AccountService) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}
