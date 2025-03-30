package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/hritesh04/shooter-game/proto"
	manager "github.com/hritesh04/shooter-game/server/gameServer/manager"
	"github.com/hritesh04/shooter-game/server/gameServer/utils"
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
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch in.Type {
		case pb.Action_Join:
			game := s.GameManger.GetRoom(in.GetRoomID())
			player := game.GetPlayer(in.GetName())
			player.AddStream(stream)
			var players []*pb.Player
			for _, p := range game.Players {
				if p.Name != in.GetName() {
					players = append(players, &pb.Player{Name: p.Name, X: float32(p.X), Y: float32(p.Y)})
				}
			}
			stream.Send(&pb.Data{Type: pb.Action_Join, Player: players})
		case pb.Action_Fire:
			game := s.GameManger.GetRoom(in.GetRoomID())
			game.EmitMove(in.GetName(), in.GetType(), in.GetData(), in.GetPlayer())
		case pb.Action_Movement:
			game := s.GameManger.GetRoom(in.GetRoomID())
			player := game.GetPlayer(in.GetName())
			player.UpdateLoc(in.GetPlayer()[0])
			game.EmitMove(in.GetName(), in.GetType(), in.GetData(), in.GetPlayer())
		}
	}
}

func (s *server) CreateRoom(ctx context.Context, data *pb.Room) (*pb.Player, error) {
	s.GameManger.AddRoom(data.GetId())
	return &pb.Player{}, nil
}

func (s *server) JoinRoom(ctx context.Context, data *pb.Room) (*pb.Room, error) {
	game := s.GameManger.GetRoom(data.GetId())
	p := game.AddPlayer()
	var players []*pb.Player
	log.Printf("Player length %d\n", len(game.Players))
	players = append(players, &pb.Player{Name: p.Name, X: float32(p.X), Y: float32(p.Y)})
	return &pb.Room{Id: game.ID, Player: players}, nil
}
func main() {
	var wg sync.WaitGroup
	host := os.Getenv("GAME_SERVER_HOST")
	if host == "" {
		// host = "shooter-local.acerowl.tech"
		host = "localhost"
	}
	port := os.Getenv("GAME_SERVER_PORT")
	if port == "" {
		port = "3001"
	}
	httpPort := os.Getenv("GAME_SERVER_HTTP_PORT")
	if httpPort == "" {
		httpPort = "3002"
	}
	grpcConn := flag.Bool("grpc", false, "Run in grpc mode")
	flag.Parse()
	wg.Add(1)
	if *grpcConn {
		log.Println("Running in grpc mode")
		go utils.RegisterServer(&wg, host+":"+port, "grpc")
	} else {
		log.Println("Running in websocket mode")
		go utils.RegisterServer(&wg, host+":"+httpPort, "websocket")
	}
	wg.Wait()
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}
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
	log.Printf("Started Websocket Server at port %s\n", httpPort)
	logger := log.New(os.Stdout, "wsproxy: ", log.LstdFlags)
	http.ListenAndServe(":"+httpPort, wsproxy.WebsocketProxy(mux, wsproxy.WithLogger(&WsLogger{logger})))
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
