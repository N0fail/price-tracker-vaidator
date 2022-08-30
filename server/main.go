package main

import (
	"context"
	"gitlab.ozon.dev/N0fail/price-tracker-validator/internal/config"
	pb "gitlab.ozon.dev/N0fail/price-tracker-validator/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conns, err := grpc.Dial(config.GrpcPortMain, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewAdminClient(conns)

	ctx = metadata.AppendToOutgoingContext(ctx, "custom", "hello")

	go runREST()
	go http.ListenAndServe("127.0.0.1:8300", nil)
	runGRPCServer(client)
}
