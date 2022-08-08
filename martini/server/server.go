package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	pb "martini/gen/proto"
	"net"
	"os"

	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/yaml.v3"
)

const serverHost = "localhost"
const serverPort = 8080
const databaseConnectionFilepath = "/server/yaml/database.yaml"

// Utils

func HandleError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Database

type DatabaseConnection struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func GetDatabaseConnection(conn DatabaseConnection, DBMS string) (*sql.DB, error) {
	conn_string := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s database=%s sslmode=disable",
		conn.Host, conn.Port, conn.User, conn.Password, conn.Database,
	)

	database, err := sql.Open(DBMS, conn_string)
	HandleError(err)

	return database, nil
}

// Server

type MartiniServer struct {
	pb.UnimplementedMartiniServer
}

func (s *MartiniServer) Echo(ctx context.Context, request *emptypb.Empty) (*pb.EchoMessage, error) {

	response := pb.EchoMessage{
		Echo: "Ping!",
	}

	return &response, nil
}

func (s *MartiniServer) GetEntity(ctx context.Context, request *pb.EntityRequest) (*pb.Entity, error) {
	// get Entity by ID from PostgreSQL

	return &pb.Entity{Id: 0, Name: "Null", Description: "Method not implemented yet"}, nil
}

// Main

func main() {

	serverAddress := fmt.Sprintf("%s:%d", serverHost, serverPort)

	wdPath, err := os.Getwd()
	HandleError(err)

	f, err := os.ReadFile(wdPath + databaseConnectionFilepath)
	HandleError(err)

	connections := make(map[string]DatabaseConnection)

	err = yaml.Unmarshal(f, &connections)
	HandleError(err)

	postgres, err := GetDatabaseConnection(connections["postgres"], "postgres")
	HandleError(err)

	defer postgres.Close()
	err = postgres.Ping()
	HandleError(err)

	fmt.Println("PostgreSQL connected!")

	listener, err := net.Listen("tcp", serverAddress)
	HandleError(err)

	grpcServer := grpc.NewServer()
	pb.RegisterMartiniServer(grpcServer, &MartiniServer{})

	err = grpcServer.Serve(listener)
	HandleError(err)

}
