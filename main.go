package main

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/tvaayi/inventory-backend/client"
	pb "github.com/tvaayi/inventory-backend/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	// waiting until db starts
	time.Sleep(10)
	mig_db, err := sql.Open("postgres", "user=postgres dbname=inventoryDB sslmode=disable password=yourpassword host=postgres_database port=5432")
	if err != nil {
		log.Fatal(err)
	}

	_, err = mig_db.Exec(`CREATE TABLE IF NOT EXISTS inventory (
		product_code TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		amount INTEGER NOT NULL,
		cost DECIMAL NOT NULL
	);`)
	if err != nil {
		log.Fatalf("Error creating inventory table: %v", err)
	}

	log.Println("Table created successfully!")

	lis, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	db, err := pgxpool.Connect(context.Background(), "postgres://postgres:yourpassword@postgres_database:5432/inventoryDB")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v", err)
	}
	s := Server{
		DB: db,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterInventoryServiceServer(grpcServer, &s)

	// run separate goruntine to run client
	go func() {
		log.Println("Starting gRPC server on port 9000...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server on port 9000: %v", err)
		}
	}()

	time.Sleep(10)
	log.Println("Waiting until server started (only for testing)")
	client.StartService()
}
