package main

import (
	"lib/internal/repository/psql"
	"lib/internal/service"
	"lib/internal/transport/rest"
	"lib/pkg/database"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	db, err := database.NewPostgresConnection(
		database.ConnectionInfo{
			Name:     "postgres",
			Port:     5432,
			Host:     "localhost",
			User:     "postgres",
			Password: "qwerty",
			SSLMode:  "disable",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	booksRepo := psql.NewBooks(db)
	service := service.NewBooks(booksRepo)
	handler := rest.NewBooksHandler(service)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
