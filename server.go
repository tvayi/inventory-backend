package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	pb "github.com/tvaayi/inventory-backend/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	DB *pgxpool.Pool
	pb.UnimplementedInventoryServiceServer
}

func (s *Server) GetInventory(ctx context.Context, req *pb.InventoryRequest) (*pb.InventoryItem, error) {
	row := s.DB.QueryRow(ctx, "SELECT product_code, name, amount, cost FROM inventory WHERE product_code=$1", req.ProductCode)

	item := pb.InventoryItem{}
	err := row.Scan(&item.ProductCode, &item.Name, &item.Amount, &item.Cost)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Inventory item not found")
	}

	return &item, nil
}

func (s *Server) AddInventoryItem(ctx context.Context, item *pb.InventoryItem) (*pb.InventoryResponse, error) {
	_, err := s.DB.Exec(ctx, "INSERT INTO inventory (product_code, name, amount, cost) VALUES ($1, $2, $3, $4)", item.ProductCode, item.Name, item.Amount, item.Cost)
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not insert inventory item")
	}

	return &pb.InventoryResponse{Response: "Inventory item successfully added"}, nil
}

func (s *Server) DeleteInventoryItem(ctx context.Context, req *pb.InventoryRequest) (*pb.InventoryResponse, error) {
	_, err := s.DB.Exec(ctx, "DELETE FROM inventory WHERE product_code=$1", req.ProductCode)
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not delete inventory item")
	}

	return &pb.InventoryResponse{Response: "Inventory item successfully deleted"}, nil
}
