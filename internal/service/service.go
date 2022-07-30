package service

import (
	"database/sql"
	"errors"
	"fmt"

	"storage/internal/entity"
)

type Interface interface {
	GetAllUsers() ([]entity.User, error)
	GetUserByID(int) (entity.User, error)
	GetUserByUsernameAndPassword(string, string) (entity.User, error)
	GetIDByUsername(string) (int, error)
	InsertUser(string, string, string) error
	DeleteUser(int) (int, error)
}

// Service ...
type Service struct {
	db *sql.DB
}

// GetService ...
func GetService(db *sql.DB) *Service {
	return &Service{db: db}
}

// GetAllUsers ...
func (s Service) GetAllUsers() (users []entity.User, err error) {
	rows, err := s.db.Query("SELECT id, username, password, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("uwu error to get all users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userBeta entity.User

		err = rows.Scan(&userBeta.ID, &userBeta.Username, &userBeta.Password, &userBeta.Email)
		if err != nil {
			return nil, fmt.Errorf("error to get all users: %w", err)
		}

		users = append(users, userBeta)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error to get all users: %w", err)
	}

	return users, nil
}

// GetUserByID ...
func (s Service) GetUserByID(id int) (user entity.User, err error) {
	row := s.db.QueryRow("SELECT id, username, password, email FROM users WHERE id = $1", id)

	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, nil
		}

		return entity.User{}, fmt.Errorf("error to get user by ID: %w", err)
	}

	return user, nil
}

// GetUserByUsernameAndPassword ...
func (s Service) GetUserByUsernameAndPassword(username, password string) (user entity.User, err error) {
	row := s.db.QueryRow(
		"SELECT id, username, password, email FROM users WHERE username = $1 AND password = $2",
		username,
		password,
	)

	err = row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, nil
		}

		return entity.User{}, fmt.Errorf("error to get user by username and password: %w", err)
	}

	return user, nil
}

// GetIDByUsername ...
func (s Service) GetIDByUsername(username string) (id int, err error) {
	row := s.db.QueryRow("SELECT id FROM users WHERE username = $1", username)

	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("error to get ID by username: %w", err)
	}

	return id, nil
}

// InsertUser ...
func (s *Service) InsertUser(username, password, email string) (err error) {
	_, err = s.db.Exec(
		"INSERT INTO users(username, password, email) VALUES ($1,$2,$3)",
		username,
		password,
		email,
	)
	if err != nil {
		return fmt.Errorf("error to insert user: %w", err)
	}

	return nil
}

// DeleteUser ...
func (s *Service) DeleteUser(id int) (rowsAffected int, err error) {
	r, err := s.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return 0, fmt.Errorf("error to delete user: %w", err)
	}

	count, _ := r.RowsAffected()

	rowsAffected = int(count)

	return rowsAffected, nil
}
