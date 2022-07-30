package service_test

import (
	"context"
	"storage/service"
	"testing"

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
			name:        nameNoError,
			outID:       idTest,
			outUsername: usernameTest,
			outEmail:    emailTest,
			inRequest:   service.EmptyRequest{},
			outErr:      "",
		},
		{
			name:        nameErrorDBClosed,
			outID:       idTest,
			outUsername: usernameTest,
			outEmail:    emailTest,
			inRequest:   service.EmptyRequest{},
			outErr:      errDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == nameErrorDBClosed {
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

			mock.ExpectQuery("^SELECT id, username, email FROM users").WillReturnRows(rows)

			r, err := service.MakeGetAllUsersEndpoint(svc)(context.TODO(), tt.inRequest)
			if err != nil {
				assert.Error(t, err)
			}

			result, ok := r.(service.UsersErrorResponse)
			if !ok {
				assert.Fail(t, "response is not of the type indicated")
			}

			if tt.name == nameNoError {
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
			name:       nameNoError,
			inID:       idTest,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			inRequest:  service.IDRequest{ID: idTest},
			outErr:     "",
		},
		{
			name: nameErrorRequest,
			inRequest: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:       nameErrorDBClosed,
			inID:       idTest,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			inRequest:  service.IDRequest{},
			outErr:     errDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			db, mock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == nameErrorDBClosed {
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

			mock.ExpectQuery("^SELECT id, username, password, email FROM users").
				WithArgs(tt.inID).WillReturnRows(rows)

			r, err := service.MakeGetUserByIDEndpoint(svc)(context.TODO(), tt.inRequest)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(service.UserErrorResponse)
			if !ok {
				if tt.name != nameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == nameNoError {
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
			name:       nameNoError,
			inID:       idTest,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			inRequest: service.UsernamePasswordRequest{
				Username: usernameTest,
				Password: passwordTest,
			},
			outErr: "",
		},
		{
			name: nameErrorRequest,
			inRequest: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:       nameErrorDBClosed,
			inID:       idTest,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			inRequest:  service.UsernamePasswordRequest{},
			outErr:     errDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			db, mock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == nameErrorDBClosed {
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

			mock.ExpectQuery("^SELECT id, username, password, email FROM users").
				WithArgs(tt.inUsername, tt.inPassword).WillReturnRows(rows)

			r, err := service.MakeGetUserByUsernameAndPasswordEndpoint(svc)(
				context.TODO(),
				tt.inRequest,
			)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(service.UserErrorResponse)
			if !ok {
				if tt.name != nameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == nameNoError {
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
			name:       nameNoError,
			inID:       idTest,
			inUsername: usernameTest,
			inRequest: service.UsernameRequest{
				Username: usernameTest,
			},
			outErr: "",
		},
		{
			name: nameErrorRequest,
			inRequest: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:       nameErrorDBClosed,
			inID:       idTest,
			inUsername: usernameTest,
			inRequest:  service.UsernameRequest{},
			outErr:     errDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			db, mock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == nameErrorDBClosed {
				db.Close()
			}

			svc := service.GetService(db)

			rows := sqlmock.NewRows([]string{"id"}).AddRow(tt.inID)

			mock.ExpectQuery("^SELECT id FROM users").WithArgs(tt.inUsername).WillReturnRows(rows)

			r, err := service.MakeGetIDByUsernameEndpoint(svc)(context.TODO(), tt.inRequest)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(service.IDErrorResponse)
			if !ok {
				if tt.name != nameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == nameNoError {
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
			name:       nameNoError,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			inRequest: service.UsernamePasswordEmailRequest{
				Username: usernameTest,
				Password: passwordTest,
				Email:    emailTest,
			},
			outErr: "",
		},
		{
			name: nameErrorRequest,
			inRequest: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:       nameErrorDBClosed,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			inRequest:  service.UsernamePasswordEmailRequest{},
			outErr:     errDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			db, mock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == nameErrorDBClosed {
				db.Close()
			}

			svc := service.GetService(db)

			mock.ExpectExec("^INSERT INTO users").
				WithArgs(
					tt.inUsername,
					tt.inPassword,
					tt.inEmail,
				).WillReturnResult(sqlmock.NewResult(0, 1))

			r, err := service.MakeInsertUserEndpoint(svc)(context.TODO(), tt.inRequest)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(service.ErrorResponse)
			if !ok {
				if tt.name != nameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == nameNoError {
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
			name: nameNoError,
			inID: idTest,
			inRequest: service.IDRequest{
				ID: idTest,
			},
			outErr: "",
		},
		{
			name: nameErrorRequest,
			inRequest: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:      nameErrorDBClosed,
			inID:      idTest,
			inRequest: service.IDRequest{},
			outErr:    errDatabaseClosed,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			db, mock, err := sqlmock.New()
			if err != nil {
				assert.Error(t, err)
			}
			defer db.Close()

			if tt.name == nameErrorDBClosed {
				db.Close()
			}

			svc := service.GetService(db)

			mock.ExpectExec("^DELETE FROM users").
				WithArgs(tt.inID).WillReturnResult(sqlmock.NewResult(0, 1))

			r, err := service.MakeDeleteUserEndpoint(svc)(context.TODO(), tt.inRequest)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(service.RowsErrorResponse)
			if !ok {
				if tt.name != nameErrorRequest {
					assert.Fail(t, "response is not of the type indicated")
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == nameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}
