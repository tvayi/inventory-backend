package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	pb "github.com/tvaayi/inventory-backend/proto"
	"google.golang.org/grpc"
)

const (
	address = ":9080"
)

func StartService() {

	// set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// initialize the InventoryService client.
	client := pb.NewInventoryServiceClient(conn)

	// call AddInventoryItem
	inventoryItem := &pb.InventoryItem{ProductCode: "PC001", Name: "Item 1", Amount: 10, Cost: 15.0}
	addRes, err := client.AddInventoryItem(context.Background(), inventoryItem)
	if err != nil {
		log.Fatalf("Could not add inventory: %v", err)
	}
	fmt.Printf("Add Inventory Response: %s\n", addRes)

	// call GetInventory
	inventoryRequest := &pb.InventoryRequest{ProductCode: "PC001"}
	res, err := client.GetInventory(context.Background(), inventoryRequest)
	if err != nil {
		log.Fatalf("Could not get inventory: %v", err)
	}
	fmt.Printf("Inventory: %s\n", res)

	// call DeleteInventoryItem
	//delRes, err := client.DeleteInventoryItem(context.Background(), inventoryRequest)
	//if err != nil {
	//	log.Fatalf("Could not delete inventory: %v", err)
	//}
	//fmt.Printf("Delete Response: %s\n", delRes)
}
