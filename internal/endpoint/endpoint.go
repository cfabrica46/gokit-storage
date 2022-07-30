package endpoint

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"storage/internal/entity"
	"storage/internal/service"

	"github.com/go-kit/kit/endpoint"
)

var ErrRequest = errors.New("error to request")

// MakeGetAllUsersEndpoint ...
func MakeGetAllUsersEndpoint(svc service.Interface) endpoint.Endpoint {
	return func(_ context.Context, _ any) (any, error) {
		var errMessage string

		users, err := svc.GetAllUsers()
		if err != nil {
			errMessage = err.Error()
		}

		return entity.UsersErrorResponse{Users: users, Err: errMessage}, nil
	}
}

// MakeGetUserByIDEndpoint ...
func MakeGetUserByIDEndpoint(svc service.Interface) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.IDRequest)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		user, err := svc.GetUserByID(req.ID)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.UserErrorResponse{User: user, Err: errMessage}, nil
	}
}

// MakeGetUserByUsernameAndPasswordEndpoint ...
func MakeGetUserByUsernameAndPasswordEndpoint(svc service.Interface) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.UsernamePasswordRequest)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		passwordHashed := NewHashHex(req.Password)

		user, err := svc.GetUserByUsernameAndPassword(req.Username, passwordHashed)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.UserErrorResponse{User: user, Err: errMessage}, nil
	}
}

// MakeGetIDByUsernameEndpoint ...
func MakeGetIDByUsernameEndpoint(svc service.Interface) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.UsernameRequest)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		id, err := svc.GetIDByUsername(req.Username)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.IDErrorResponse{ID: id, Err: errMessage}, nil
	}
}

// MakeInsertUserEndpoint ...
func MakeInsertUserEndpoint(svc service.Interface) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.UsernamePasswordEmailRequest)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		passwordHashed := NewHashHex(req.Password)

		err := svc.InsertUser(req.Username, passwordHashed, req.Email)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.ErrorResponse{Err: errMessage}, nil
	}
}

// MakeDeleteUserEndpoint ...
func MakeDeleteUserEndpoint(svc service.Interface) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.IDRequest)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		rowsAffected, err := svc.DeleteUser(req.ID)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.RowsErrorResponse{RowsAffected: rowsAffected, Err: errMessage}, nil
	}
}

func NewHashHex(data string) (hash string) {
	hasher := sha256.New()

	hasher.Write([]byte(data))

	return hex.EncodeToString(hasher.Sum(nil))
}
