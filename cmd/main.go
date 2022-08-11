package main

import (
	"database/sql"
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
	log.SetFlags(log.Lshortfile)

	if !config.VerifyIsDockerRun() {
		if err := config.LoadEnv(); err != nil {
			log.Fatal(err)
		}
	}

	dbInfo := config.GetDBInfo()

	db, err := sql.Open(config.DBDriver, dbInfo)
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
