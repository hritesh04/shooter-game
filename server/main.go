package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	pb "github.com/hritesh04/shooter-game/stubs"
	"google.golang.org/grpc"
)

type server struct {
	GameManger GameManager
	pb.UnimplementedMovementEmitterServer
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

type Game struct {
	ID      string
	Players []*Player
	Started bool
}

func (g *Game) AddPlayer() *Player {
	var player *Player
	if len(g.Players) == 0 {
		player = &Player{
			Name: generateSecureID(),
			X:    60,
			Y:    70,
		}
	} else {
		player = &Player{
			Name: generateSecureID(),
			X:    1172,
			Y:    608,
		}
		for _, p := range g.Players {
			log.Printf("Player %s data sent to %s\n", player.Name, p.Name)
			p.Conn.Send(&pb.Data{Type: pb.Action_Info, Player: []*pb.Player{&pb.Player{Name: player.Name, X: float32(player.X), Y: float32(player.Y)}}})
		}
	}
	g.Players = append(g.Players, player)
	log.Printf("Player %s joined room %s\n", player.Name, g.ID)
	return player
}

func (g *Game) GetPlayer(name string) *Player {
	var player *Player
	for _, p := range g.Players {
		if p.Name == name {
			player = p
		}
	}
	return player
}

func (g *Game) EmitMove(name string, direction pb.Direction) {
	data := &pb.Data{
		Type: pb.Action_Movement,
		Data: direction,
	}
	for _, p := range g.Players {
		if p.Name == name {
			continue
		}
		p.Conn.Send(data)
	}
}

// var game map[string]Game

type GameManager struct {
	games map[string]*Game
}

func (g *GameManager) AddRoom(roomID string) {
	_, ok := g.games[roomID]
	if ok {
		// gameID = generateSecureID()
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
		// gameID = generateSecureID()
		log.Printf("dungeon not present %s", roomID)
		return nil
	}
	// game.AddPlayer(roomID)
	return game
	// g.games[roomID] = game
}

func generateSecureID() string {
	b := make([]byte, 3) // 3 bytes = 6 hex characters
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *server) SendMove(stream pb.MovementEmitter_SendMoveServer) error {
	// if type join get the roomID and player name and add the stream to the Game.player.stream and return game and player info
	for {
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			fmt.Printf("Receiving message from stream: %v\n", err)
			return err
		}
		fmt.Printf("bidi echoing message %q\n", in)
		switch in.Type {
		case pb.Action_Join:
			game := s.GameManger.GetRoom(in.GetRoomID())
			player := game.GetPlayer(in.GetName())
			player.AddStream(stream)
			// player, err := getPlayerInfo(in.GetRoomID(), in.GetName())
			// if err != nil {
			// 	return err
			// }
			// player.Conn = stream
			log.Printf("player joined dungeon with stream : %+v\n", player)
			// stream.Send(&pb.Data{Type: pb.Action_Join, Player: []*pb.Player{&pb.Player{Name: player.Name, X: float32(player.X), Y: float32(player.Y)}}})
		case pb.Action_Fire:
		case pb.Action_Movement:
		}
		// if in.Type == pb.Action_Join {
		// 	roomID := in.GetRoomID()
		// 	fmt.Printf("Player %s joined Room %s", in.GetName(), roomID)
		// 	players := game[roomID].Players
		// 	var player Player
		// 	for _, p := range players {
		// 		if p.Name == in.GetName() {
		// 			player = p
		// 		}
		// 	}
		// 	player.Conn = stream
		// 	log.Printf("Sent data %+v", player)
		// 	stream.Send(&pb.Data{Type: pb.Action_Info, Player: &pb.Player{Name: player.Name, X: float32(player.X), Y: float32(player.Y)}})
		// }

	}
}

func (s *server) CreateRoom(ctx context.Context, data *pb.Room) (*pb.Player, error) {
	// gameID := data.GetId()
	// game :=
	s.GameManger.AddRoom(data.GetId())
	// gameID = generateSecureID()
	// _, ok := game[gameID]
	// _, ok := game[gameID]
	// if ok {
	// gameID = generateSecureID()
	// return nil, fmt.Errorf("dungeon already present")
	// }
	// player := Player{
	// 	Name: data.GetName(),
	// 	X:    60,
	// 	Y:    70,
	// }
	// g := Game{
	// 	ID:      gameID,
	// 	Started: false,
	// }
	// g.Players = append(g.Players, player)
	// game[gameID] = g
	// if g.Started {
	// 	return nil, fmt.Errorf("dungeon battle stated")
	// }
	// log.Printf("created room : %s\n", gameID)
	// if err != nil {
	// 	if err == io.EOF {
	// 		return nil
	// 	}
	// fmt.Printf("Receiving message from stream: %v\n", err)
	// fmt.Printf("bidi echoing message %q\n", in.Move)
	// stream.Send(&pb.Data{Move: in.Move, Direction: in.Direction})
	return &pb.Player{}, nil
}

func (s *server) JoinRoom(ctx context.Context, data *pb.Room) (*pb.Room, error) {
	// fmt.Println("%+v", data)
	game := s.GameManger.GetRoom(data.GetId())
	player := game.AddPlayer()
	// fmt.Printf("Player %s joined room %s\n", player.Name, data.GetId())
	var players []*pb.Player
	for _, p := range game.Players {
		players = append(players, &pb.Player{Name: p.Name, X: float32(p.X), Y: float32(p.Y)})
	}
	return &pb.Room{Id: game.ID, Name: player.Name, Player: players}, nil
}

// func getPlayerInfo(roomID string, name string) (*Player, error) {
// 	game, ok := game[roomID]
// 	if !ok {
// 		return nil, fmt.Errorf("dungeon not found")
// 	}
// 	players := game.Players
// 	var player Player
// 	for _, p := range players {
// 		if p.Name == name {
// 			player = p
// 		}
// 	}
// 	return &player, nil
// }

func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}
	go func() {
		server := struct {
			Address string `json:"address"`
		}{
			Address: "localhost:3000",
		}
		out, err := json.Marshal(server)
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/registerServer", bytes.NewBuffer(out))
		if err != nil {
			log.Printf("error creating request : %w", err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("error making http call : %w", err)
		}
		if res.StatusCode != http.StatusOK {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Printf("error parsing response body :%w", err)
			}
			log.Fatalf("error in registrating server :%s", string(body))
		}
	}()
	// game = make(map[string]Game)
	s := grpc.NewServer()
	pb.RegisterMovementEmitterServer(s, &server{
		GameManger: GameManager{games: make(map[string]*Game)},
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
