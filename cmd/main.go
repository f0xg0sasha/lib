package main

import (
	"fmt"
	"lib/internal/config"
	"lib/internal/repository/psql"
	"lib/internal/service"
	grpc_client "lib/internal/transport/grpc"
	"lib/internal/transport/rest"
	"lib/pkg/database"
	"lib/pkg/hash"
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
	hasher := hash.NewSHA1Hasher("salt")

	auditService, err := grpc_client.NewClient(9000)
	if err != nil {
		log.Fatal(err)
	}

	booksRepo := psql.NewBooks(db)
	booksService := service.NewBooks(booksRepo, auditService)

	usersRepo := psql.NewUsers(db)
	tokenRepo := psql.NewToken(db)

	usersService := service.NewUsers(usersRepo, tokenRepo, hasher, auditService, []byte(os.Getenv("HASH_SECRET")), cfg.Auth.TokenTTL)

	handler := rest.NewHandler(booksService, usersService)

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
