package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"database/sql"
	// "github.com/jackc/pgx"
	pb "../part"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Part) (*pb.Part, error)
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Part
}

// Create a new consignment
func (repo *Repository) Create(consignment *pb.Part) (*pb.Part, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

type PartServiceServer struct {
}

func (s *PartServiceServer) CreateParts(ctx context.Context, req *pb.CreatePartReq) (*pb.CreatePartRes, error) {
	// Essentially doing req.Blog to access the struct with a nil check
	part := req.GetPart()
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// defer conn.Close(context.Background())
	defer conn.Close()

	// for part:= range parts {
	err = conn.QueryRow(
		"INSERT INTO part (id, vendor_code, manufacturer_id) VALUES (DEFAULT, $1, $2) RETURNING id",
		part.VendorCode, part.ManufacturerId).Scan(&part.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	// }
	return &pb.CreatePartRes{Part: part}, nil
}

func (s *PartServiceServer) ReadPart(ctx context.Context, req *pb.ReadPartReq) (*pb.ReadPartRes, error) {
	partID := req.GetId()
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	part := &pb.Part{Id: partID}
	err = conn.QueryRow(
		"SELECT vendor_code, manufacturer_id FROM part where id=$1",
		partID).Scan(&part.VendorCode, &part.ManufacturerId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		// os.Exit(1)
		part = nil
	}
	return &pb.ReadPartRes{Part: part}, nil
}

func (s *PartServiceServer) UpdatePart(ctx context.Context, req *pb.UpdatePartReq) (*pb.UpdatePartRes, error) {
	part := req.GetPart()
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Exec("UPDATE part SET vendor_code=$2, manufacturer_id=$3 where id=$1", part.Id, part.VendorCode, part.ManufacturerId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return &pb.UpdatePartRes{Part: nil}, nil
	}
	return &pb.UpdatePartRes{Part: part}, nil
}

func (s *PartServiceServer) DeletePart(ctx context.Context, req *pb.DeletePartReq) (*pb.DeletePartRes, error) {
	partID := req.GetId()
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Exec("UPDATE part SET deleted_at=NOW() where id=$1", partID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		return &pb.DeletePartRes{Success: false}, nil
	}
	return &pb.DeletePartRes{Success: true}, nil
}

func main() {
	m, err := migrate.New(
		"file://./migrations",
		os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterPartCrudServiceServer(s, &PartServiceServer{})

	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
