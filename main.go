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
}

func (s *server) SendNotification(ctx context.Context, in *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	fmt.Println("Sending notification:", in.body)
	for _, conn := range s.wsConnections {
		conn.WriteJSON(map[string]string{
			"title":   "Notification Title",
			"message": in.body,
			"icon":    "path/to/icon.png",
		})
	}
	return &pb.NotificationResponse{message: "Notification sent!"}, nil
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
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// HTTP server
	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
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
