package service_test

import (
	"testing"

	"storage/internal/entity/mock"
	"storage/internal/service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name                         string
		outID, outUsername, outEmail any
		outErr                       string
	}{
		{
			name:        mock.NameNoError,
			outID:       mock.IDTest,
			outUsername: mock.UsernameTest,
			outEmail:    mock.EmailTest,
			outErr:      "",
		},
		{
			name:        mock.NameErrorDBClosed,
			outID:       mock.IDTest,
			outUsername: mock.UsernameTest,
			outEmail:    mock.EmailTest,
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

			dbMock.ExpectQuery("SELECT id, username, email FROM users").WillReturnRows(rows)

			_, err = svc.GetAllUsers()
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

func TestGetUserByID(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name                               string
		outUsername, outPassword, outEmail string
		outErr                             string
		inID                               int
	}{
		{
			name:        mock.NameNoError,
			inID:        mock.IDTest,
			outUsername: mock.UsernameTest,
			outPassword: mock.PasswordTest,
			outEmail:    mock.EmailTest,
			outErr:      "",
		},
		{
			name:        mock.NameErrorNoRows,
			inID:        mock.IDTest,
			outUsername: mock.UsernameTest,
			outPassword: mock.PasswordTest,
			outEmail:    mock.EmailTest,
			outErr:      "",
		},
		{
			name:        mock.NameErrorDBClosed,
			inID:        mock.IDTest,
			outUsername: mock.UsernameTest,
			outPassword: mock.PasswordTest,
			outEmail:    mock.EmailTest,
			outErr:      "sql: database is closed",
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
				tt.outUsername,
				tt.outPassword,
				tt.outEmail,
			)

			if tt.name == mock.NameErrorNoRows {
				rows = sqlmock.NewRows([]string{"id", "username", "password", "email"})
			}

			dbMock.ExpectQuery(
				"^SELECT id, username, password, email FROM users",
			).WithArgs(tt.inID).WillReturnRows(rows)

			_, err = svc.GetUserByID(tt.inID)
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

func TestGetUserByUsernameAndPassword(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
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
			outErr:     "",
		},
		{
			name:       mock.NameErrorNoRows,
			inID:       mock.IDTest,
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			outErr:     "",
		},
		{
			name:       mock.NameErrorDBClosed,
			inID:       mock.IDTest,
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			outErr:     "sql: database is closed",
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

			if tt.name == mock.NameErrorNoRows {
				rows = sqlmock.NewRows([]string{"id", "username", "password", "email"})
			}

			dbMock.ExpectQuery(
				"^SELECT id, username, password, email FROM users",
			).WithArgs(tt.inUsername, tt.inPassword).WillReturnRows(rows)

			user, err := svc.GetUserByUsernameAndPassword(tt.inUsername, tt.inPassword)
			if err != nil {
				resultErr = err.Error()
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, resultErr)
				assert.NotEmpty(t, user)
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
			name:       mock.NameNoError,
			inID:       mock.IDTest,
			inUsername: mock.UsernameTest,
			outErr:     "",
		},
		{
			name:       mock.NameErrorNoRows,
			inID:       mock.IDTest,
			inUsername: mock.UsernameTest,
			outErr:     "",
		},
		{
			name:       mock.NameErrorDBClosed,
			inID:       mock.IDTest,
			inUsername: mock.UsernameTest,
			outErr:     "sql: database is closed",
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

			if tt.name == mock.NameErrorNoRows {
				rows = sqlmock.NewRows([]string{"id"})
			}

			dbMock.ExpectQuery("^SELECT id FROM users").WithArgs(tt.inUsername).WillReturnRows(rows)

			_, err = svc.GetIDByUsername(tt.inUsername)
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

func TestInsertUser(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name                            string
		inUsername, inPassword, inEmail string
		outErr                          string
	}{
		{
			name:       mock.NameNoError,
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			outErr:     "",
		},
		{
			name:       mock.NameErrorNoRows,
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			outErr:     "",
		},
		{
			name:       mock.NameErrorDBClosed,
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			outErr:     "sql: database is closed",
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

			dbMock.ExpectExec(
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

			if tt.name == mock.NameNoError {
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
			name:   mock.NameNoError,
			inID:   mock.IDTest,
			outErr: "",
		},
		{
			name:   mock.NameErrorDBClosed,
			inID:   mock.IDTest,
			outErr: "sql: database is closed",
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

			dbMock.ExpectExec(
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

			if tt.name == mock.NameNoError {
				assert.Empty(t, resultErr)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}
