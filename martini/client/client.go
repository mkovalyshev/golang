package main

import (
	"context"
	"fmt"
	"log"
	pb "martini/gen/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}

	client := pb.NewMartiniClient(conn)

	response, err := client.Echo(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Println(err)
	}

	fmt.Println(response.Echo)
}
