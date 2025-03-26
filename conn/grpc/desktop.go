package grpc

import (
	"context"
	"io"
	"log"

	pb "github.com/hritesh04/shooter-game/proto"
	"github.com/hritesh04/shooter-game/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DesktopClient struct {
	address string
	client  pb.MovementEmitterClient
	conn    pb.MovementEmitter_SendMoveClient
}

func NewGrpcDesktopClient(address string) types.IConnection {

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error creating coonnection,%v", err)
	}

	client := pb.NewMovementEmitterClient(conn)

	log.Println("Connected using GRPC")
	return &DesktopClient{
		address: address,
		client:  client,
	}
}

func (c *DesktopClient) JoinRoom(ID string) (*pb.Room, error) {
	ctx := context.Background()
	join, err := c.client.JoinRoom(ctx, &pb.Room{Id: ID})
	if err != nil {
		log.Fatalf("error joining room %v", err)
		return nil, err
	}
	return join, nil
}

func (c *DesktopClient) createEventConn() {
	ctx := context.Background()
	conn, err := c.client.SendMove(ctx)
	if err != nil {
		log.Println(err)
	}
	c.conn = conn
}

func (c *DesktopClient) GetEventConn() pb.MovementEmitter_SendMoveClient {
	if c.conn == nil {
		c.createEventConn()
	}
	return c.conn
}

func (c *DesktopClient) SendMove(conn grpc.BidiStreamingClient[pb.Data, pb.Data], data *pb.Data) error {
	log.Printf("sending data from move :%v", data)
	err := conn.Send(data)
	if err == io.EOF {
		return err
	}
	if err != nil {
		log.Fatalf("sending Streaming message: %v", err)
		return err
	}
	closeErr := conn.CloseSend()
	if closeErr != nil {
		log.Fatalf("cannot close send: %w", err)
		return err
	}
	return nil
}
