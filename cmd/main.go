package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"storage/cmd/config"
	"storage/internal/endpoint"
	"storage/internal/entity"
	"storage/internal/service"
	"storage/internal/transport"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.GetAPIConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := openPostgresConn(cfg.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	runServer(cfg.Port, db)
}

func runServer(port string, db *sql.DB) {
	svc := service.GetService(db)

	getAllUsersHandler := httptransport.NewServer(
		endpoint.MakeGetAllUsersEndpoint(svc),
		transport.DecodeRequestWithoutBody(),
		transport.EncodeResponse,
	)

	getUserByIDHandler := httptransport.NewServer(
		endpoint.MakeGetUserByIDEndpoint(svc),
		transport.DecodeRequest(entity.IDRequest{}),
		transport.EncodeResponse,
	)

	getUserByUsernameAndPasswordHandler := httptransport.NewServer(
		endpoint.MakeGetUserByUsernameAndPasswordEndpoint(svc),
		transport.DecodeRequest(entity.UsernamePasswordRequest{}),
		transport.EncodeResponse,
	)

	getIDByUsernameHandler := httptransport.NewServer(
		endpoint.MakeGetIDByUsernameEndpoint(svc),
		transport.DecodeRequest(entity.UsernameRequest{}),
		transport.EncodeResponse,
	)

	insertUserHandler := httptransport.NewServer(
		endpoint.MakeInsertUserEndpoint(svc),
		transport.DecodeRequest(entity.UsernamePasswordEmailRequest{}),
		transport.EncodeResponse,
	)

	deleteUserHandler := httptransport.NewServer(
		endpoint.MakeDeleteUserEndpoint(svc),
		transport.DecodeRequest(entity.IDRequest{}),
		transport.EncodeResponse,
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

func openPostgresConn(conn config.DBConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conn.Host,
		conn.Port,
		conn.User,
		conn.Password,
		conn.DBName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
