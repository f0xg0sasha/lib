package main

import (
	"fmt"
	"lib/internal/config"
	"lib/internal/repository/psql"
	"lib/internal/service"
	"lib/internal/transport/rest"
	"lib/pkg/database"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	cfg, err := config.NewConfig(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading.env file: %s", err.Error())
	}

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

	log.Info("SERVER STARTED AT")

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
