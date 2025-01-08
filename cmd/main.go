package main

import (
	"fmt"
	"lib/internal/config"
	"lib/internal/repository/psql"
	"lib/internal/service"
	"lib/internal/transport/rest"
	"lib/pkg/database"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func main() {
	cfg, err := config.NewConfig(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading.env file: %s", err.Error())
	}

	fmt.Println(cfg)

	db, err := database.NewPostgresConnection(
		database.ConnectionInfo{
			Name:     os.Getenv("DB_DBNAME"),
			Port:     StringToInt(os.Getenv("DB_PORT")),
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_NAME"),
			Password: os.Getenv("DB_PASSWORD"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
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
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
