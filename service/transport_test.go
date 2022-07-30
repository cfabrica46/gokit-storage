package service_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"storage/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myRequests struct {
	idReq, usernameReq, usernamePasswordReq, usernamePasswordEmailReq, badReq *http.Request
}

const (
	idRequestJSON = `{
		 "id": 1
	}`

	usernameRequestJSON = `{
		 "username": "username"
	}`

	//nolint:gosec
	usernamePasswordRequestJSON = `{
		 "username": "username",
		 "password": "password"
	}`

	//nolint:gosec
	usernamePasswordEmailRequestJSON = `{
		 "username": "username",
		 "password": "password",
		 "email": "email@email.com"
	}`
)

func TestDecodeRequestWithoutBody(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		out    any
		in     *http.Request
		name   string
		outErr string
	}{
		{
			name:   nameNoError + "IDRequest",
			in:     nil,
			outErr: "",
			out:    service.EmptyRequest{},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := service.DecodeRequestWithoutBody()(context.TODO(), tt.in)

			assert.Empty(t, err)
			assert.Equal(t, tt.out, r)
		})
	}
}

func TestDecodeRequest(t *testing.T) {
	t.Parallel()

	myReqs, err := getRequests()
	if err != nil {
		assert.Error(t, err)
	}

	for _, tt := range []struct {
		inType      any
		in          *http.Request
		name        string
		outUsername string
		outPassword string
		outEmail    string
		outErr      string
		outID       int
	}{
		{
			name:   nameNoError + "IDRequest",
			inType: service.IDRequest{},
			in:     myReqs.idReq,
			outID:  idTest,
			outErr: "",
		},
		{
			name:        nameNoError + "UsernameRequest",
			inType:      service.UsernameRequest{},
			in:          myReqs.usernameReq,
			outUsername: usernameTest,
			outErr:      "",
		},
		{
			name:        nameNoError + "UsernamePasswordRequest",
			inType:      service.UsernamePasswordRequest{},
			in:          myReqs.usernamePasswordReq,
			outUsername: usernameTest,
			outPassword: passwordTest,
			outErr:      "",
		},
		{
			name:        nameNoError + "UsernamePasswordEmailRequest",
			inType:      service.UsernamePasswordEmailRequest{},
			in:          myReqs.usernamePasswordEmailReq,
			outUsername: usernameTest,
			outPassword: passwordTest,
			outEmail:    emailTest,
			outErr:      "",
		},
		{
			name:   "BadRequest",
			inType: service.IDRequest{},
			in:     myReqs.badReq,
			outErr: "EOF",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			var req any

			switch resultType := tt.inType.(type) {
			case service.IDRequest:
				req, err = service.DecodeRequest(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(service.IDRequest)
				if ok {
					assert.Equal(t, tt.outID, result.ID)
					assert.Contains(t, resultErr, tt.outErr)
				} else {
					assert.NotNil(t, err)
				}

			case service.UsernameRequest:
				req, err = service.DecodeRequest(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(service.UsernameRequest)
				assert.True(t, ok)

				assert.Equal(t, tt.outUsername, result.Username)
				assert.Contains(t, resultErr, tt.outErr)

			case service.UsernamePasswordRequest:
				req, err = service.DecodeRequest(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(service.UsernamePasswordRequest)
				assert.True(t, ok)

				assert.Equal(t, tt.outUsername, result.Username)
				assert.Equal(t, tt.outPassword, result.Password)
				assert.Contains(t, resultErr, tt.outErr)

			case service.UsernamePasswordEmailRequest:
				req, err = service.DecodeRequest(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(service.UsernamePasswordEmailRequest)
				assert.True(t, ok)

				assert.Equal(t, tt.outUsername, result.Username)
				assert.Equal(t, tt.outPassword, result.Password)
				assert.Equal(t, tt.outEmail, result.Email)
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}

func TestEncodeResponse(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		in     any
		outErr string
	}{
		{
			name:   nameNoError,
			in:     "test",
			outErr: "",
		},
		{
			name:   "ErrorBadEncode",
			in:     func() {},
			outErr: "json: unsupported type: func()",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var resultErr string

			err := service.EncodeResponse(context.TODO(), httptest.NewRecorder(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			if tt.name == nameNoError {
				assert.Empty(t, resultErr)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}

func getRequests() (myReqs *myRequests, err error) {
	idReq, err := http.NewRequest(
		http.MethodPost,
		urlTest,
		bytes.NewBuffer([]byte(idRequestJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	usernameReq, err := http.NewRequest(
		http.MethodPost,
		urlTest,
		bytes.NewBuffer([]byte(usernameRequestJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	usernamePasswordReq, err := http.NewRequest(
		http.MethodPost,
		urlTest,
		bytes.NewBuffer([]byte(usernamePasswordRequestJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	usernamePasswordEmailReq, err := http.NewRequest(
		http.MethodPost,
		urlTest,
		bytes.NewBuffer([]byte(usernamePasswordEmailRequestJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	badReq, err := http.NewRequest(http.MethodPost, urlTest, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return &myRequests{
		idReq:                    idReq,
		usernameReq:              usernameReq,
		usernamePasswordReq:      usernamePasswordReq,
		usernamePasswordEmailReq: usernamePasswordEmailReq,
		badReq:                   badReq,
	}, nil
}
