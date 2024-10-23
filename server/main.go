package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/hritesh04/shooter-game/stubs"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMovementEmitterServer
}

func (s *server) SendMove(stream pb.MovementEmitter_SendMoveServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Printf("Receiving message from stream: %v\n", err)
			return err
		}
		fmt.Printf("bidi echoing message %q\n", in.Move)
		stream.Send(&pb.Movement{Move: in.Move})
	}
}

func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMovementEmitterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
