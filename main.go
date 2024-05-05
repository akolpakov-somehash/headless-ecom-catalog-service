package main

import (
	"catalog/internal"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	pb "github.com/akolpakov-somehash/go-microservices/proto/catalog/product"
	"github.com/cenkalti/backoff/v4"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"google.golang.org/grpc"
)

func loadEnv() error {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}

func connectDB() (*gorm.DB, error) {
	// Get environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	var db *gorm.DB
	var err error

	operation := func() error {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		return err
	}

	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.MaxElapsedTime = 2 * time.Minute
	exponentialBackOff.MaxInterval = 10 * time.Second

	err = backoff.Retry(operation, exponentialBackOff)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

func startServer(db *gorm.DB, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &internal.Server{DB: db})
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}

const (
	defaultPort = 50051
)

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	db, err := connectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&internal.DbProduct{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	port := defaultPort
	if p, ok := os.LookupEnv("PORT"); ok {
		port, err = strconv.Atoi(p)
		if err != nil {
			log.Fatalf("invalid port number: %v", err)
		}
	}

	err = startServer(db, port)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
