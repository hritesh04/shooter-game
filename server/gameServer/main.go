package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	pb "github.com/hritesh04/shooter-game/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type server struct {
	GameManger GameManager
	pb.UnimplementedMovementEmitterServer
}
type Game struct {
	ID      string
	Players []Player
	Started bool
}

func (g *Game) AddPlayer() *Player {
	var player Player
	if len(g.Players) == 0 {
		player = Player{
			Name: generateSecureID(),
			X:    60,
			Y:    70,
		}
	} else {
		player = Player{
			Name: generateSecureID(),
			X:    1172,
			Y:    608,
		}
		fmt.Printf("%+v", g.Players)
		for _, p := range g.Players {
			log.Printf("Player %s data sent to %s\n", player.Name, p.Name)
			fmt.Println("HERE")
			p.Conn.Send(&pb.Data{Type: pb.Action_Info, Player: []*pb.Player{&pb.Player{Name: player.Name, X: float32(player.X), Y: float32(player.Y)}}})
		}
	}
	g.Players = append(g.Players, player)
	log.Printf("Player %s joined room %s\n", player.Name, g.ID)
	return &player
}

func (g *Game) GetPlayer(name string) *Player {
	for i := range g.Players {
		if g.Players[i].Name == name {
			return &g.Players[i]
		}
	}
	return nil
}

func (g *Game) EmitMove(name string, action pb.Action, direction pb.Direction, player []*pb.Player) {
	data := &pb.Data{
		Type:   action,
		Data:   direction,
		Name:   name,
		Player: player,
	}
	for _, p := range g.Players {
		if p.Name == name {
			continue
		}
		p.Conn.Send(data)
	}
}

type Player struct {
	Name string
	Conn pb.MovementEmitter_SendMoveServer
	X    float64
	Y    float64
}

func (p *Player) AddStream(stream pb.MovementEmitter_SendMoveServer) {
	p.Conn = stream
}

func (p *Player) UpdateLoc(player *pb.Player) {
	p.X = float64(player.GetX())
	p.Y = float64(player.GetY())
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
	fmt.Printf("ctx %+v", ctx)
	game := s.GameManger.GetRoom(data.GetId())
	p := game.AddPlayer()
	var players []*pb.Player
	log.Printf("Player length %d", len(game.Players))
	players = append(players, &pb.Player{Name: p.Name, X: float32(p.X), Y: float32(p.Y)})
	return &pb.Room{Id: game.ID, Player: players}, nil
}
func main() {
	host := os.Getenv("GAME_SERVER_HOST")
	if host == "" {
		host = "shooter-local.acerowl.tech"
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
		}{
			Address: host,
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
	// game = make(map[string]Game)
	opts := []grpc.ServerOption{grpc.UnaryInterceptor(corsInterceptor), grpc.MaxRecvMsgSize(1024 * 1024 * 64), grpc.MaxSendMsgSize(1024 * 1024 * 64)}
	s := grpc.NewServer(opts...)
	pb.RegisterMovementEmitterServer(s, &server{
		GameManger: GameManager{games: make(map[string]*Game)},
	})
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr))
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
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

type GameManager struct {
	games map[string]*Game
}

func (g *GameManager) AddRoom(roomID string) {
	_, ok := g.games[roomID]
	if ok {
		log.Printf("dungeon already present %s", roomID)
		return
	}
	game := &Game{
		ID:      roomID,
		Started: false,
	}
	g.games[roomID] = game
	log.Printf("created room : %s\n", roomID)
}

func (g *GameManager) GetRoom(roomID string) *Game {
	game, ok := g.games[roomID]
	if !ok {
		log.Printf("dungeon not present %s", roomID)
		return nil
	}
	return game
}

func generateSecureID() string {
	b := make([]byte, 3) // 3 bytes = 6 hex characters
	rand.Read(b)
	return hex.EncodeToString(b)
}
