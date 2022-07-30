package service_test

import (
	"storage/service"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	urlTest string = "localhost:8080"

	idTest       int    = 1
	usernameTest string = "username"
	passwordTest string = "password"
	emailTest    string = "email@email.com"

	errDatabaseClosed string = "sql: database is closed"

	nameNoError       string = "NoError"
	nameErrorRequest  string = "ErrorRequest"
	nameErrorDBClosed string = "ErrorDBClosed"
	nameErrorNoRows   string = "ErrorNoRows"
)

func TestGetAllUsers(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name                         string
		outID, outUsername, outEmail any
		outErr                       string
	}{
		{
			name:        nameNoError,
			outID:       idTest,
			outUsername: usernameTest,
			outEmail:    emailTest,
			outErr:      "",
		},
		{
			name:        nameErrorDBClosed,
			outID:       idTest,
			outUsername: usernameTest,
			outEmail:    emailTest,
			outErr:      "sql: database is closed",
		},
		{
			name:        "ErrorScanRows",
			outID:       "id",
			outUsername: 1,
			outEmail:    1,
			outErr:      "Scan error on column index 0",
		},
		/* {
			name:        nameErrorNoRows,
			outID:       idTest,
			outUsername: usernameTest,
			outEmail:    emailTest,
			outErr:      "asdfadfsafds",
		}, */
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
					"email",
				}).AddRow(
				tt.outID,
				tt.outUsername,
				tt.outEmail,
			)

			mock.ExpectQuery("SELECT id, username, email FROM users").WillReturnRows(rows)

			_, err = svc.GetAllUsers()
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

func TestGetUserByID(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
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
			outErr:     "",
		},
		{
			name:       nameErrorNoRows,
			inID:       idTest,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			outErr:     "",
		},
		{
			name:       nameErrorDBClosed,
			inID:       idTest,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			outErr:     "sql: database is closed",
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

			if tt.name == nameErrorNoRows {
				rows = sqlmock.NewRows([]string{"id", "username", "password", "email"})
			}

			mock.ExpectQuery(
				"^SELECT id, username, password, email FROM users",
			).WithArgs(tt.inID).WillReturnRows(rows)

			_, err = svc.GetUserByID(tt.inID)
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

func TestGetUserByUsernameAndPassword(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
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
			outErr:     "",
		},
		{
			name:       nameErrorNoRows,
			inID:       idTest,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			outErr:     "",
		},
		{
			name:       nameErrorDBClosed,
			inID:       idTest,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			outErr:     "sql: database is closed",
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

			if tt.name == nameErrorNoRows {
				rows = sqlmock.NewRows([]string{"id", "username", "password", "email"})
			}

			mock.ExpectQuery(
				"^SELECT id, username, password, email FROM users",
			).WithArgs(tt.inUsername, tt.inPassword).WillReturnRows(rows)

			_, err = svc.GetUserByUsernameAndPassword(tt.inUsername, tt.inPassword)
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

func TestGetIDByUsername(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name       string
		inUsername string
		outErr     string
		inID       int
	}{
		{
			name:       nameNoError,
			inID:       idTest,
			inUsername: usernameTest,
			outErr:     "",
		},
		{
			name:       nameErrorNoRows,
			inID:       idTest,
			inUsername: usernameTest,
			outErr:     "",
		},
		{
			name:       nameErrorDBClosed,
			inID:       idTest,
			inUsername: usernameTest,
			outErr:     "sql: database is closed",
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

			if tt.name == nameErrorNoRows {
				rows = sqlmock.NewRows([]string{"id"})
			}

			mock.ExpectQuery("^SELECT id FROM users").WithArgs(tt.inUsername).WillReturnRows(rows)

			_, err = svc.GetIDByUsername(tt.inUsername)
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

func TestInsertUser(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name                            string
		inUsername, inPassword, inEmail string
		outErr                          string
	}{
		{
			name:       nameNoError,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			outErr:     "",
		},
		{
			name:       nameErrorNoRows,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			outErr:     "",
		},
		{
			name:       nameErrorDBClosed,
			inUsername: usernameTest,
			inPassword: passwordTest,
			inEmail:    emailTest,
			outErr:     "sql: database is closed",
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

			mock.ExpectExec(
				"^INSERT INTO users",
			).WithArgs(
				tt.inUsername,
				tt.inPassword,
				tt.inEmail,
			).WillReturnResult(
				sqlmock.NewResult(0, 1),
			)

			err = svc.InsertUser(tt.inUsername, tt.inPassword, tt.inEmail)
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

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		outErr string
		inID   int
	}{
		{
			name:   nameNoError,
			inID:   idTest,
			outErr: "",
		},
		{
			name:   nameErrorDBClosed,
			inID:   idTest,
			outErr: "sql: database is closed",
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

			mock.ExpectExec(
				"^DELETE FROM users",
			).WithArgs(
				tt.inID,
			).WillReturnResult(
				sqlmock.NewResult(0, 1),
			)

			_, err = svc.DeleteUser(tt.inID)
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
