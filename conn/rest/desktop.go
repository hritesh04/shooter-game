package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	pb "github.com/hritesh04/shooter-game/proto"
	"github.com/hritesh04/shooter-game/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type DesktopClient struct {
	address string
	conn    pb.MovementEmitter_SendMoveClient
}

type Connection struct {
	conn *websocket.Conn
}

type WebSocketenvelope struct {
	Result WebSocketMessage `json:"result"`
}

type WebSocketMessage struct {
	Type   string       `json:"type"`
	Data   string       `json:"data"`
	Name   string       `json:"name,omitempty"`
	RoomID string       `json:"roomID,omitempty"`
	Player []*pb.Player `json:"player,omitempty"`
}

type Player struct {
	Name string  `json:"name"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

func (c Connection) Send(data *pb.Data) error {
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling req %v", err)
		return err
	}
	if err := c.conn.WriteMessage(websocket.BinaryMessage, payload); err != nil {
		log.Println("Error sending movement data through websocket %v", err)
		return err
	}
	return nil
}

func (c Connection) Recv() (*pb.Data, error) {
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		log.Println("Error receiving movement data from websocket %v", err)
	}
	var wsMessage WebSocketenvelope
	if err := json.Unmarshal(msg, &wsMessage); err != nil {
		log.Printf("Error unmarshalling websocket message: %v", err)
		return nil, err
	}
	pbData := pb.Data{
		Type: convertStringToAction(wsMessage.Result.Type),
		Data: convertStringToDirection(wsMessage.Result.Data),
	}
	players := make([]*pb.Player, len(wsMessage.Result.Player))
	for i, p := range wsMessage.Result.Player {
		players[i] = &pb.Player{Name: p.Name, X: float32(p.X), Y: float32(p.Y)}
	}
	pbData.Player = players
	return &pbData, nil
}

func convertStringToAction(actionStr string) pb.Action {
	switch actionStr {
	case "Join":
		return pb.Action_Join
	case "Movement":
		return pb.Action_Movement
	case "Fire":
		return pb.Action_Fire
	case "Info":
		return pb.Action_Info
	default:
		return pb.Action_Info
	}
}

func convertStringToDirection(directionStr string) pb.Direction {
	switch directionStr {
	case "UP":
		return pb.Direction_UP
	case "DOWN":
		return pb.Direction_DOWN
	case "LEFT":
		return pb.Direction_LEFT
	case "RIGHT":
		return pb.Direction_RIGHT
	default:
		return pb.Direction_UP
	}
}

func (c Connection) CloseSend() error {
	return nil
}
func (c Connection) Context() context.Context {
	return nil
}

func (c Connection) Header() (metadata.MD, error) {
	return nil, nil
}

func (c Connection) RecvMsg(m any) error {
	return nil
}

func (c Connection) SendMsg(m any) error {
	return nil
}
func (c Connection) Trailer() metadata.MD {
	return nil
}

func NewRestDesktopClient(address string) types.IConnection {
	log.Println("Connected using Websocket")
	return &DesktopClient{
		address: address,
	}
}

func (c *DesktopClient) JoinRoom(ID string) (*pb.Room, error) {
	data, err := json.Marshal(&pb.Room{Id: ID})
	if err != nil {
		log.Println("Error marshaling req %v", err)
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, "http://"+c.address+"/v1/joinRoom", bytes.NewReader(data))
	if err != nil {
		log.Println("Error creating JoinRoom Http Request %v", err)
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read request: %v\n", err)
		return nil, err
	}
	res := &pb.Room{}
	if err := json.Unmarshal(body, res); err != nil {
		log.Printf("failed to unmarshal response: %v\n", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("Error respone from join room %s", resp.Status)
		return nil, fmt.Errorf("error: received non 200 status code")
	}
	return res, nil
}

func (c *DesktopClient) createEventConn() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://"+c.address+"/v1/sendMove", nil)
	if err != nil {
		log.Println("Failed to connect to WebSocket server: %v", err)
	}
	c.conn = &Connection{
		conn: conn,
	}
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
