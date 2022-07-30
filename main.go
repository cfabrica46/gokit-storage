package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"storage/service"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	log.SetFlags(log.Lshortfile)

	if err := godotenv.Load(".env"); err != nil {
		log.Println(err)
	}

	dbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	fmt.Println(dbInfo)

	db, err := sql.Open(os.Getenv("DB_DRIVER"), dbInfo)
	if err != nil {
		log.Println(err)

		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)

		return
	}

	runServer(os.Getenv("PORT"), db)
}

func runServer(port string, db *sql.DB) {
	svc := service.GetService(db)

	getAllUsersHandler := httptransport.NewServer(
		service.MakeGetAllUsersEndpoint(svc),
		service.DecodeRequestWithoutBody(),
		service.EncodeResponse,
	)

	getUserByIDHandler := httptransport.NewServer(
		service.MakeGetUserByIDEndpoint(svc),
		service.DecodeRequest(service.IDRequest{}),
		service.EncodeResponse,
	)

	getUserByUsernameAndPasswordHandler := httptransport.NewServer(
		service.MakeGetUserByUsernameAndPasswordEndpoint(svc),
		service.DecodeRequest(service.UsernamePasswordRequest{}),
		service.EncodeResponse,
	)

	getIDByUsernameHandler := httptransport.NewServer(
		service.MakeGetIDByUsernameEndpoint(svc),
		service.DecodeRequest(service.UsernameRequest{}),
		service.EncodeResponse,
	)

	insertUserHandler := httptransport.NewServer(
		service.MakeInsertUserEndpoint(svc),
		service.DecodeRequest(service.UsernamePasswordEmailRequest{}),
		service.EncodeResponse,
	)

	deleteUserHandler := httptransport.NewServer(
		service.MakeDeleteUserEndpoint(svc),
		service.DecodeRequest(service.IDRequest{}),
		service.EncodeResponse,
	)

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/users").Handler(getAllUsersHandler)
	router.Methods(http.MethodGet).Path("/user/id").Handler(getUserByIDHandler)
	router.Methods(http.MethodGet).Path("/user/username_password").
		Handler(getUserByUsernameAndPasswordHandler)
	router.Methods(http.MethodGet).Path("/id/username").Handler(getIDByUsernameHandler)
	router.Methods(http.MethodPost).Path("/user").Handler(insertUserHandler)
	router.Methods(http.MethodDelete).Path("/user").Handler(deleteUserHandler)

	log.Println("ListenAndServe on localhost:" + os.Getenv("PORT"))
	log.Println(http.ListenAndServe(":"+port, router))
}
