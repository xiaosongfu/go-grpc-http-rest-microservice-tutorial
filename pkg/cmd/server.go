package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"go.xiaosongfu.com/go-grpc-http-rest-microservice-tutorial/pkg/protocol/rest"

	_ "github.com/go-sql-driver/mysql"
	"go.xiaosongfu.com/go-grpc-http-rest-microservice-tutorial/pkg/protocol/grpc"
	"go.xiaosongfu.com/go-grpc-http-rest-microservice-tutorial/pkg/service/v1"
)

type Config struct {
	GRPCPort            string
	HttpPort            string
	DatastoreDBHost     string
	DatastoreDBPort     string
	DatastoreDBUser     string
	DatastoreDBPassword string
	DatastoreDBSchema   string
}

func RunServer() error {
	ctx := context.Background()

	var cfg Config

	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.HttpPort, "http-port", "", "http port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DatastoreDBPort, "db-port", "", "Database port")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")

	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}
	if len(cfg.HttpPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfg.HttpPort)
	}

	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBPort,
		cfg.DatastoreDBSchema,
		param)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	defer db.Close()

	v1API := v1.NewToDoServiceServer(db)

	// 启动 http gateway
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HttpPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
