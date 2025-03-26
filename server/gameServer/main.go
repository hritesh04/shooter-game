package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/hritesh04/shooter-game/proto"
	manager "github.com/hritesh04/shooter-game/server/gameServer/manager"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type server struct {
	GameManger *manager.GameManager
	pb.UnimplementedMovementEmitterServer
}

func (s *server) SendMove(stream pb.MovementEmitter_SendMoveServer) error {
	for {
		in, err := stream.Recv()
		fmt.Println(in)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch in.Type {
		case pb.Action_Join:
			fmt.Println("Join Action by client")
			game := s.GameManger.GetRoom(in.GetRoomID())
			player := game.GetPlayer(in.GetName())
			player.AddStream(stream)
			var players []*pb.Player
			for _, p := range game.Players {
				if p.Name != in.GetName() {
					players = append(players, &pb.Player{Name: p.Name, X: float32(p.X), Y: float32(p.Y)})
				}
			}
			fmt.Println("Join Action sent by server")
			stream.Send(&pb.Data{Type: pb.Action_Join, Player: players})
		case pb.Action_Fire:
			fmt.Println("Fire Action by client")
			game := s.GameManger.GetRoom(in.GetRoomID())
			game.EmitMove(in.GetName(), in.GetType(), in.GetData(), in.GetPlayer())
		case pb.Action_Movement:
			fmt.Println("Movement Action by client")
			game := s.GameManger.GetRoom(in.GetRoomID())
			player := game.GetPlayer(in.GetName())
			player.UpdateLoc(in.GetPlayer()[0])
			fmt.Println("Movement Action sent by server")
			game.EmitMove(in.GetName(), in.GetType(), in.GetData(), in.GetPlayer())
		}
	}
}

func (s *server) CreateRoom(ctx context.Context, data *pb.Room) (*pb.Player, error) {
	s.GameManger.AddRoom(data.GetId())
	fmt.Println("Room Creation func triggered")
	fmt.Println(data.GetId())
	return &pb.Player{}, nil
}

func (s *server) JoinRoom(ctx context.Context, data *pb.Room) (*pb.Room, error) {
	fmt.Println("join Room func triggered")
	game := s.GameManger.GetRoom(data.GetId())
	p := game.AddPlayer()
	var players []*pb.Player
	log.Printf("Player length %d", len(game.Players))
	players = append(players, &pb.Player{Name: p.Name, X: float32(p.X), Y: float32(p.Y)})
	fmt.Println(data.GetId())
	return &pb.Room{Id: game.ID, Player: players}, nil
}
func main() {
	host := os.Getenv("GAME_SERVER_HOST")
	if host == "" {
		// host = "shooter-local.acerowl.tech"
		host = "localhost:3002"
	}
	port := os.Getenv("GAME_SERVER_PORT")
	if port == "" {
		port = "3001"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}
	go func() {
		server := struct {
			Address string `json:"address"`
			Type    string `json:"type"`
		}{
			Address: host,
			Type:    "rest",
		}
		out, err := json.Marshal(server)
		if err != nil {
			log.Fatalf("error binding request payload : %v", err)
		}
		req, err := http.NewRequest(http.MethodPost, os.Getenv("ROOT_SERVER_URL")+"/registerServer", bytes.NewBuffer(out))
		if err != nil {
			log.Printf("error creating request : %v", err)
		}
		tlsVerification := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tlsVerification}

		res, err := client.Do(req)
		if err != nil {
			log.Printf("error making http call : %v", err)
		}
		if res != nil {
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				body, err := io.ReadAll(res.Body)
				if err != nil {
					log.Printf("error parsing response body :%v", err)
					return
				}
				log.Printf("error in registrating server :%s", string(body))
				return
			}
			log.Printf("Successfully registered server with root server")
		}
	}()
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(corsInterceptor),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterMovementEmitterServer(s, &server{
		GameManger: manager.NewGameManager(),
	})
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve : %v", err)
		}
		log.Printf("server listening at %v", lis.Addr())
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	restOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	restErr := pb.RegisterMovementEmitterHandlerFromEndpoint(ctx, mux, "localhost:"+port, restOpts)
	if restErr != nil {
		log.Fatalf(restErr.Error())
	}
	fmt.Println("listening")
	logger := log.New(os.Stdout, "wsproxy: ", log.LstdFlags)
	http.ListenAndServe(":3002", wsproxy.WebsocketProxy(mux, wsproxy.WithLogger(&WsLogger{logger})))
}

type WsLogger struct {
	*log.Logger
}

func (w *WsLogger) Debugln(v ...interface{}) {
	w.Println("[DEBUG]", v)
}

func (w *WsLogger) Warnln(v ...interface{}) {
	w.Println("[WARN]", v)
}

func corsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	headers := metadata.Pairs(
		"Access-Control-Allow-Origin", "*",
		"Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE",
		"Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization",
	)
	grpc.SetHeader(ctx, headers)
	return handler(ctx, req)
}
