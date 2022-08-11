package config

import (
	"fmt"
	"os"
)

const (
	DBDriver = "postgres"
)

func VerifyIsDockerRun() (check bool) {
	isDocker := os.Getenv("DOCKER")

	return isDocker != ""
}

func GetDBInfo() (dbInfo string) {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
}

func LoadEnv() (err error) {
	err = os.Setenv("PORT", "7070")
	if err != nil {
		return fmt.Errorf("error to set env:%w", err)
	}

	err = os.Setenv("DB_HOST", "localhost")
	if err != nil {
		return fmt.Errorf("error to set env:%w", err)
	}

	err = os.Setenv("DB_PORT", "5431")
	if err != nil {
		return fmt.Errorf("error to set env:%w", err)
	}

	err = os.Setenv("DB_USERNAME", "cfabrica46")
	if err != nil {
		return fmt.Errorf("error to set env:%w", err)
	}

	err = os.Setenv("DB_PASSWORD", "01234")
	if err != nil {
		return fmt.Errorf("error to set env:%w", err)
	}

	err = os.Setenv("DB_NAME", "go_crud")
	if err != nil {
		return fmt.Errorf("error to set env:%w", err)
	}

	err = os.Setenv("DB_SSLMODE", "disable")
	if err != nil {
		return fmt.Errorf("error to set env:%w", err)
	}

	return nil
}
