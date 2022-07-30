package transport_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"storage/internal/entity"
	"storage/internal/entity/mock"
	"storage/internal/transport"

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
			name:   mock.NameNoError + "IDRequest",
			in:     nil,
			outErr: "",
			out:    entity.EmptyRequest{},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := transport.DecodeRequestWithoutBody()(context.TODO(), tt.in)

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
			name:   mock.NameNoError + "IDRequest",
			inType: entity.IDRequest{},
			in:     myReqs.idReq,
			outID:  mock.IDTest,
			outErr: "",
		},
		{
			name:        mock.NameNoError + "UsernameRequest",
			inType:      entity.UsernameRequest{},
			in:          myReqs.usernameReq,
			outUsername: mock.UsernameTest,
			outErr:      "",
		},
		{
			name:        mock.NameNoError + "UsernamePasswordRequest",
			inType:      entity.UsernamePasswordRequest{},
			in:          myReqs.usernamePasswordReq,
			outUsername: mock.UsernameTest,
			outPassword: mock.PasswordTest,
			outErr:      "",
		},
		{
			name:        mock.NameNoError + "UsernamePasswordEmailRequest",
			inType:      entity.UsernamePasswordEmailRequest{},
			in:          myReqs.usernamePasswordEmailReq,
			outUsername: mock.UsernameTest,
			outPassword: mock.PasswordTest,
			outEmail:    mock.EmailTest,
			outErr:      "",
		},
		{
			name:   "BadRequest",
			inType: entity.IDRequest{},
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
			case entity.IDRequest:
				req, err = transport.DecodeRequest(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(entity.IDRequest)
				if ok {
					assert.Equal(t, tt.outID, result.ID)
					assert.Contains(t, resultErr, tt.outErr)
				} else {
					assert.NotNil(t, err)
				}

			case entity.UsernameRequest:
				req, err = transport.DecodeRequest(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(entity.UsernameRequest)
				assert.True(t, ok)

				assert.Equal(t, tt.outUsername, result.Username)
				assert.Contains(t, resultErr, tt.outErr)

			case entity.UsernamePasswordRequest:
				req, err = transport.DecodeRequest(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(entity.UsernamePasswordRequest)
				assert.True(t, ok)

				assert.Equal(t, tt.outUsername, result.Username)
				assert.Equal(t, tt.outPassword, result.Password)
				assert.Contains(t, resultErr, tt.outErr)

			case entity.UsernamePasswordEmailRequest:
				req, err = transport.DecodeRequest(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(entity.UsernamePasswordEmailRequest)
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
			name:   mock.NameNoError,
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

			err := transport.EncodeResponse(context.TODO(), httptest.NewRecorder(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			if tt.name == mock.NameNoError {
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
		mock.URLTest,
		bytes.NewBuffer([]byte(idRequestJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	usernameReq, err := http.NewRequest(
		http.MethodPost,
		mock.URLTest,
		bytes.NewBuffer([]byte(usernameRequestJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	usernamePasswordReq, err := http.NewRequest(
		http.MethodPost,
		mock.URLTest,
		bytes.NewBuffer([]byte(usernamePasswordRequestJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	usernamePasswordEmailReq, err := http.NewRequest(
		http.MethodPost,
		mock.URLTest,
		bytes.NewBuffer([]byte(usernamePasswordEmailRequestJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	badReq, err := http.NewRequest(http.MethodPost, mock.URLTest, bytes.NewBuffer([]byte{}))
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
