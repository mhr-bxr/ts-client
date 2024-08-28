package main

import (
	"context"
	"log"
	"time"

	pb "github.com/mhr-bxr/ts-proto/tsapi"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "server_address:port", grpc.WithInsecure(), grpc.WithBlock()) //nolint:staticcheck
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTonServiceClient(conn)

	// Create a transaction request
	txData := []byte{0x00, 0x01, 0x02} // example transaction data
	req := &pb.TransactionRequest{
		Tx:            txData,
		WalletVersion: "1.0",
	}

	// Call the ProcessTx method
	res, err := c.ProcessTx(context.Background(), req)
	if err != nil {
		log.Fatalf("could not process transaction: %v", err)
	}
	log.Printf("Transaction Processed: UUID = %s", res.Value)
}
