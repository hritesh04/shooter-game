package common

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/hritesh04/shooter-game/stubs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SendMove(ctx context.Context, client pb.MovementEmitterClient) error {
	c, err := client.SendMove(ctx)
	if err != nil {
		log.Fatalf("error make function call %v", err)
	}
	go func() {
		for {
			// resp, err := c.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("receiving Streaming message: %v", err)
			}
			// fmt.Println("BidiStreaming Msg: ", resp.Move)
		}
	}()
	for i := 0; i < 100; i++ {
		// err := c.Send(&pb.Data{Move: pb.Action_Movement, Direction: pb.Direction_LEFT})
		if err == io.EOF {
			// Bidi streaming RPC errors happen and make Send return io.EOF,
			// not the RPC error itself.  Call Recv to determine the error.
			break
		}
		if err != nil {
			// Some local errors are reported this way, e.g. errors serializing
			// the request message.
			log.Fatalf("sending Streaming message: %v", err)
		}
	}
	err = c.CloseSend()
	if err != nil {
		return fmt.Errorf("cannot close send: %v", err)
	}
	return nil
}

func NewGrpcClient(address string) pb.MovementEmitterClient {

	// to make sure that go doesnt exited after one run in wasm
	// <-make(chan struct{})

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error creating coonnection,%v", err)
	}
	// defer conn.Close()
	// fmt.Println(conn)
	client := pb.NewMovementEmitterClient(conn)
	// ctx := context.Background()
	// defer cancel()
	// if err := SendMove(ctx, client); err != nil {
	// 	log.Fatalf("Bidirectional Streaming RPC failed: %v", err)
	// }
	// time.Sleep(time.Second * 10)
	return client
}
