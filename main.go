package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	pb "github.com/shishir248/Notification/notifications"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Define ports for connections
const (
	grpcPort = ":50051"
	httpPort = ":8080"
)

// Get websocket upgrader
var upgrader = websocket.Upgrader{}

type server struct {
	wsConnections []*websocket.Conn
	pb.UnimplementedNotificationServer
}

func (s *server) Send(ctx context.Context, in *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	fmt.Println("Sending notification:", in.Body)
	for _, conn := range s.wsConnections {
		conn.WriteJSON(map[string]string{
			"title":   "Notification Title",
			"message": in.Body,
			"icon":    "path/to/icon.png",
		})
	}
	return &pb.NotificationResponse{Message: "Notification sent!"}, nil
}

func main() {
	s := &server{}

	// gRPC server
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServer(grpcServer, s)
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		} else {
			fmt.Println("Server is running")
		}
	}()
	// HTTP server
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		s.wsConnections = append(s.wsConnections, conn)
		defer conn.Close()
	})
	http.ListenAndServe(httpPort, r)
}
