package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "path/to/proto/file"
)

type server struct{}

func (s *server) Subscribe(ctx context.Context, in *pb.Subscription) (*pb.Response, error) {
	return &pb.Response{Message: "Subscribed Successfully"}, nil
}

func (s *server) Unsubscribe(ctx context.Context, in *pb.Subscription) (*pb.Response, error) {
	return &pb.Response{Message: "Unsubscribed Successfully"}, nil
}

func (s *server) SendNotification(ctx context.Context, in *pb.Notification) (*pb.Response, error) {
	// Send the notification to subscribed clients
	return &pb.Response{Message: "Notification sent"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPushNotificationServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
