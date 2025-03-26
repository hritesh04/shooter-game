package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"

	pb "github.com/hritesh04/shooter-game/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PlayerRequest struct {
	RoomID string `json:"roomID"`
}

type PlayerResponse struct {
	RoomID  string `json:"roomID"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

type ErrorResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"err"`
}

type Server struct {
	Address          string   `json:"address"`
	RoomID           []string `json:"roomID"`
	Type             string   `json:"type"`
	ActiveConnection int      `json:"activeConnection"`
	MaxConnection    int      `json:"maxConnection"`
}

func (s *Server) AddRoom(roomID string) error {
	if s.Type == "grpc" {
		conn, err := grpc.NewClient(s.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("error creating coonnection,%v", err)
			return errors.New("error connecting to grpc game server")
		}
		client := pb.NewMovementEmitterClient(conn)
		ctx := context.Background()
		_, err = client.CreateRoom(ctx, &pb.Room{Id: roomID})
		if err != nil {
			log.Printf("error creating room at remote server : %s", s.Address)
			log.Println(err)
			return errors.New("error creating room at grpc gamer server")
		}
	} else {
		data, err := json.Marshal(&pb.Room{Id: roomID})
		fmt.Println(data)
		if err != nil {
			log.Printf("Error marshaling req %v", err)
			return errors.New("error marshaling request")
		}
		req, err := http.NewRequest(http.MethodPost, "http://"+s.Address+"/v1/createRoom", bytes.NewReader(data))
		if err != nil {
			log.Println("error creating request", err)
			return errors.New("error creating request")
		}
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Println("error making request", err)
			return errors.New("error making request to game server")
		}
		defer res.Body.Close()
	}
	s.RoomID = append(s.RoomID, roomID)
	return nil
}

type ServerManager struct {
	severMap map[string]Server
	Server   []Server
}

func (s *ServerManager) GetLeastConnectedServer() Server {
	leastConn := math.MaxInt
	var server Server
	for _, s := range s.Server {
		if s.ActiveConnection < leastConn {
			server = s
		}
	}
	return server
}

func (s *ServerManager) AddServer(server Server) bool {
	for _, s := range s.Server {
		if s.Address == server.Address {
			return true
		}
	}
	s.Server = append(s.Server, server)
	return true
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	port := os.Getenv("ROOT_SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	serverManager := &ServerManager{
		severMap: make(map[string]Server),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/createRoom", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var player PlayerRequest
		if err := decoder.Decode(&player); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Invalid Input"))
		}
		server := serverManager.GetLeastConnectedServer()
		if player.RoomID == "" {
			player.RoomID = generateSecureID()
		}
		if _, ok := serverManager.severMap[player.RoomID]; ok {
			player.RoomID = generateSecureID()
		}
		if err := server.AddRoom(player.RoomID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		serverManager.severMap[player.RoomID] = server
		response := PlayerResponse{
			Address: server.Address,
			Type:    server.Type,
			RoomID:  player.RoomID,
		}
		log.Printf("New room created: %s at server: %s", response.RoomID, response.Address)
		w.WriteHeader(http.StatusOK)
		writeJSON(w, response)
	})
	mux.HandleFunc("/joinRoom", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var player PlayerRequest
		if err := decoder.Decode(&player); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		server, ok := serverManager.severMap[player.RoomID]
		if !ok {
			http.Error(w, fmt.Errorf("dungeon not found").Error(), http.StatusBadRequest)
			return
		}
		server.ActiveConnection++
		response := PlayerResponse{
			Address: server.Address,
			RoomID:  player.RoomID,
		}
		log.Printf("New join room request current count %d", server.ActiveConnection)
		log.Printf("New player joined room: %s at server: %s", response.RoomID, response.Address)
		w.WriteHeader(http.StatusOK)
		writeJSON(w, response)
	})
	mux.HandleFunc("/registerServer", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Errorf("invalid method").Error(), http.StatusMethodNotAllowed)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var server Server
		if err := decoder.Decode(&server); err != nil {
			http.Error(w, fmt.Errorf("failed to decode body").Error(), http.StatusBadRequest)
			return
		}
		if flag := serverManager.AddServer(server); !flag {
			log.Printf("Server already registered url: %s", server.Address)
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Printf("New server registered url: %s", server.Address)
		w.WriteHeader(http.StatusOK)
	})
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	log.Printf("Starting Server listening at port %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed")
	}
}

func generateSecureID() string {
	b := make([]byte, 3) // 3 bytes = 6 hex characters
	rand.Read(b)
	return hex.EncodeToString(b)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}
}
