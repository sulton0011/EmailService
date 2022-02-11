package main

import (
	"fmt"
	"net"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/home_work/TaskUserService/EmailService/config"
	"github.com/home_work/TaskUserService/EmailService/pkg/logger"
	"github.com/home_work/TaskUserService/EmailService/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/home_work/TaskUserService/EmailService/genproto/email_service"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "email_service")
	defer logger.Cleanup(log)

	log.Info("main: pgxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase)

	connDb := sqlx.MustConnect("postgres", psqlString)

	sendService := service.NewSendService(connDb, cfg)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("error listening tcp port: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterEmailServiceServer(s, sendService)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
