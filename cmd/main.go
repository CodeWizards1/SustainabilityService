package main

import (
	"fmt"
	"log"
	"net"
	"sustainabilityService/config"
	pb "sustainabilityService/genproto/SustainabilityService"
	"sustainabilityService/services"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func GetDB(path string) (*sqlx.DB, error) {
	cfg := config.Load(path)

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.DbHost,
		cfg.Postgres.DbPort,
		cfg.Postgres.DbUser,
		cfg.Postgres.DbPassword,
		cfg.Postgres.DbName,
	)

	db, err := sqlx.Connect("postgres", psqlUrl)
	return db, err
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50053")
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	gprcServer := grpc.NewServer()

	db, err := GetDB(".")
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	sustainabilityService, err := services.NewSustainabilityService(db, "50052")
	if err != nil {
		log.Fatal("Failed to create sustainability service: ", err)
	}
	pb.RegisterSustainabilityServiceServer(gprcServer, sustainabilityService)

	log.Println("gRPC server is running on port 50053")
	if err := gprcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
