package main

import (
	"context"
	"github.com/fatih/color"
	"github.com/tlb_katia/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	address = "localhost:50051"
	noteID  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := auth_v1.NewAuthV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &auth_v1.GetRequest{Id: noteID})
	if err != nil {
		log.Fatalf("failed to get auth by id: %v", err)
	}

	log.Printf(color.RedString("Client info:\n"), color.GreenString("%+v", r.GetName()))
}
