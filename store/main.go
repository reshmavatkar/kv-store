package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "github.com/reshmavatkar/kv-store/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StoreServe is interface for store
type StoreServe interface {
	Put(ctx context.Context, in *pb.PutRequest) (*pb.PutResponse, error)
	Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error)
	Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error)
}

type storeServer struct {
	pb.UnimplementedKeyValueStoreServer
	store map[string]string
	mu    sync.RWMutex
}

// NewStoreServer creates new instance for store server
func NewStoreServer() pb.KeyValueStoreServer {
	return &storeServer{
		store: make(map[string]string),
	}
}

// Put stores Key Value in in-memory
func (s *storeServer) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key is required")
	}
	if len(req.Value) == 0 {
		return nil, status.Error(codes.InvalidArgument, "value cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[req.Key] = req.Value
	return &pb.PutResponse{Status: "OK"}, nil
}

// Get retrives value for the key
func (s *storeServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key is required")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.store[req.Key]
	if !ok {
		return nil, status.Error(codes.NotFound, "key not found")
	}
	return &pb.GetResponse{Value: val}, nil
}

// Delete deletes Key value from in-memory
func (s *storeServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.store[req.Key]
	if !exists {
		return nil, status.Error(codes.NotFound, "key not found")
	}
	delete(s.store, req.Key)
	return &pb.DeleteResponse{Success: true}, nil
}

func main() {
	grpcPort := ":50051"
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	storeSvc := NewStoreServer()
	pb.RegisterKeyValueStoreServer(grpcServer, storeSvc)
	// TODO: add UnaryInterceptor Logging, auth, metrics, panic recovery for unary RPCs

	go func() {
		log.Printf("gRPC Store service is listening on %s\n", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down gRPC server...")

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("gRPC server stopped gracefully")
	case <-time.After(10 * time.Second):
		log.Println("gRPC server shutdown timed out, forcing stop")
		grpcServer.Stop()
	}
}
