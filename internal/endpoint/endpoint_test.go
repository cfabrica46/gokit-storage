package endpoint_test

import (
	"context"
	"testing"

	"storage/internal/endpoint"
	"storage/internal/entity"
	"storage/internal/entity/mock"
	"storage/internal/service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type incorrectRequest struct {
	incorrect bool
}

func TestMakeGetAllUsersEndpoint(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		inRequest                          any
		name                               string
		outUsername, outPassword, outEmail string
		outErr                             string
		outID                              int
	}{
		{
			name:        mock.NameNoError,
			outID:       mock.IDTest,
			outUsername: mock.UsernameTest,
			outEmail:    mock.EmailTest,
			inRequest:   entity.EmptyRequest{},
			outErr:      "",
		},
		{
			name:        mock.NameErrorDBClosed,
			outID:       mock.IDTest,
			outUsername: mock.UsernameTest,
			outEmail:    mock.EmailTest,
			inRequest:   entity.EmptyRequest{},
			outErr:      mock.ErrDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db, dbMock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == mock.NameErrorDBClosed {
				db.Close()
			}

			svc := service.GetService(db)

			rows := sqlmock.NewRows(
				[]string{
					"id",
					"username",
					"email",
				}).AddRow(
				tt.outID,
				tt.outUsername,
				tt.outEmail,
			)

			dbMock.ExpectQuery("^SELECT id, username, email FROM users").WillReturnRows(rows)

			r, err := endpoint.MakeGetAllUsersEndpoint(svc)(context.TODO(), tt.inRequest)
			if err != nil {
				assert.Error(t, err)
			}

			result, ok := r.(entity.UsersErrorResponse)
			if !ok {
				assert.Fail(t, "response is not of the type indicated")
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, result.Err, tt.outErr)
			}
		})
	}
}

func TestMakeGetUserByIDEndpoint(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		inRequest                       any
		name                            string
		inUsername, inPassword, inEmail string
		outErr                          string
		inID                            int
	}{
		{
			name:       mock.NameNoError,
			inID:       mock.IDTest,
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			inRequest:  entity.IDRequest{ID: mock.IDTest},
			outErr:     "",
		},
		{
			name: mock.NameErrorRequest,
			inRequest: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:       mock.NameErrorDBClosed,
			inID:       mock.IDTest,
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			inRequest:  entity.IDRequest{},
			outErr:     mock.ErrDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			db, dbMock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == mock.NameErrorDBClosed {
				db.Close()
			}

			svc := service.GetService(db)

			rows := sqlmock.NewRows(
				[]string{
					"id",
					"username",
					"password",
					"email",
				}).AddRow(
				tt.inID,
				tt.inUsername,
				tt.inPassword,
				tt.inEmail,
			)

			dbMock.ExpectQuery("^SELECT id, username, password, email FROM users").
				WithArgs(tt.inID).WillReturnRows(rows)

			r, err := endpoint.MakeGetUserByIDEndpoint(svc)(context.TODO(), tt.inRequest)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.UserErrorResponse)
			if !ok {
				if tt.name != mock.NameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}

func TestMakeGetUserByUsernameAndPasswordEndpoint(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		inRequest                       any
		name                            string
		inUsername, inPassword, inEmail string
		outErr                          string
		inID                            int
	}{
		{
			name:       mock.NameNoError,
			inID:       mock.IDTest,
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			inRequest: entity.UsernamePasswordRequest{
				Username: mock.UsernameTest,
				Password: mock.PasswordTest,
			},
			outErr: "",
		},
		{
			name: mock.NameErrorRequest,
			inRequest: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:       mock.NameErrorDBClosed,
			inID:       mock.IDTest,
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			inRequest:  entity.UsernamePasswordRequest{},
			outErr:     mock.ErrDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			passwordHashed := endpoint.NewHashHex(tt.inPassword)

			db, dbMock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == mock.NameErrorDBClosed {
				db.Close()
			}

			svc := service.GetService(db)

			rows := sqlmock.NewRows(
				[]string{
					"id",
					"username",
					"password",
					"email",
				}).AddRow(
				tt.inID,
				tt.inUsername,
				passwordHashed,
				tt.inEmail,
			)

			dbMock.ExpectQuery("^SELECT id, username, password, email FROM users").
				WithArgs(tt.inUsername, passwordHashed).WillReturnRows(rows)

			r, err := endpoint.MakeGetUserByUsernameAndPasswordEndpoint(svc)(
				context.TODO(),
				tt.inRequest,
			)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.UserErrorResponse)
			if !ok {
				if tt.name != mock.NameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}

func TestGetIDByUsernameEndpoint(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		inRequest  any
		name       string
		inUsername string
		outErr     string
		inID       int
	}{
		{
			name:       mock.NameNoError,
			inID:       mock.IDTest,
			inUsername: mock.UsernameTest,
			inRequest: entity.UsernameRequest{
				Username: mock.UsernameTest,
			},
			outErr: "",
		},
		{
			name: mock.NameErrorRequest,
			inRequest: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:       mock.NameErrorDBClosed,
			inID:       mock.IDTest,
			inUsername: mock.UsernameTest,
			inRequest:  entity.UsernameRequest{},
			outErr:     mock.ErrDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			db, dbMock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == mock.NameErrorDBClosed {
				db.Close()
			}

			svc := service.GetService(db)

			rows := sqlmock.NewRows([]string{"id"}).AddRow(tt.inID)

			dbMock.ExpectQuery("^SELECT id FROM users").WithArgs(tt.inUsername).WillReturnRows(rows)

			r, err := endpoint.MakeGetIDByUsernameEndpoint(svc)(context.TODO(), tt.inRequest)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.IDErrorResponse)
			if !ok {
				if tt.name != mock.NameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}

func TestMakeInsertUserEndpoint(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		inRequest                       any
		name                            string
		inUsername, inPassword, inEmail string
		outErr                          string
	}{
		{
			name:       mock.NameNoError,
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			inRequest: entity.UsernamePasswordEmailRequest{
				Username: mock.UsernameTest,
				Password: mock.PasswordTest,
				Email:    mock.EmailTest,
			},
			outErr: "",
		},
		{
			name: mock.NameErrorRequest,
			inRequest: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:       mock.NameErrorDBClosed,
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			inRequest:  entity.UsernamePasswordEmailRequest{},
			outErr:     mock.ErrDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			passwordHashed := endpoint.NewHashHex(tt.inPassword)

			db, dbMock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == mock.NameErrorDBClosed {
				db.Close()
			}

			svc := service.GetService(db)

			dbMock.ExpectExec("^INSERT INTO users").
				WithArgs(
					tt.inUsername,
					passwordHashed,
					tt.inEmail,
				).WillReturnResult(sqlmock.NewResult(0, 1))

			r, err := endpoint.MakeInsertUserEndpoint(svc)(context.TODO(), tt.inRequest)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.ErrorResponse)
			if !ok {
				if tt.name != mock.NameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}

func TestMakeDeleteUserEndpoint(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		inRequest any
		name      string
		outErr    string
		inID      int
	}{
		{
			name: mock.NameNoError,
			inID: mock.IDTest,
			inRequest: entity.IDRequest{
				ID: mock.IDTest,
			},
			outErr: "",
		},
		{
			name: mock.NameErrorRequest,
			inRequest: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:      mock.NameErrorDBClosed,
			inID:      mock.IDTest,
			inRequest: entity.IDRequest{},
			outErr:    mock.ErrDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			db, dbMock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == mock.NameErrorDBClosed {
				db.Close()
			}

			svc := service.GetService(db)

			dbMock.ExpectExec("^DELETE FROM users").
				WithArgs(tt.inID).WillReturnResult(sqlmock.NewResult(0, 1))

			r, err := endpoint.MakeDeleteUserEndpoint(svc)(context.TODO(), tt.inRequest)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.RowsErrorResponse)
			if !ok {
				if tt.name != mock.NameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}
