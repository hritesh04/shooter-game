package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	pb "github.com/hritesh04/shooter-game/stubs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PlayerRequest struct {
	RoomID string `json:"roomID"`
}

type PlayerResponse struct {
	RoomID  string `json:"roomID"`
	Address string `json:"address"`
}

type ErrorResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"err"`
}

type Server struct {
	Address          string   `json:"address"`
	RoomID           []string `json:"roomID"`
	ActiveConnection int      `json:"activeConnection"`
	MaxConnection    int      `json:"maxConnection"`
}

func (s *Server) AddRoom(roomID string) {
	conn, err := grpc.NewClient(s.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error creating coonnection,%v", err)
		return
	}
	client := pb.NewMovementEmitterClient(conn)
	ctx := context.Background()
	_, err = client.CreateRoom(ctx, &pb.Room{Id: roomID})
	if err != nil {
		log.Printf("error creating room at remote server : %s", s.Address)
	}
	s.RoomID = append(s.RoomID, roomID)
}

type ServerManager struct {
	severMap map[string]string
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

func (s *ServerManager) AddServer(server Server) {
	s.Server = append(s.Server, server)
}

func main() {
	serverManager := &ServerManager{
		severMap: make(map[string]string),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/createRoom", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// decoder := json.NewDecoder(r.Body)
		// var player PlayerRequest
		// if err := decoder.Decode(&player); err != nil {
		// 	w.WriteHeader(http.StatusBadGateway)
		// 	w.Write([]byte("Invalid Input"))
		// }
		server := serverManager.GetLeastConnectedServer()
		var roomID string
		roomID = generateSecureID()
		if _, ok := serverManager.severMap[roomID]; !ok {
			roomID = generateSecureID()
		}
		server.AddRoom(roomID)
		serverManager.severMap[roomID] = server.Address
		// w.WriteHeader(http.StatusOK)
		response := PlayerResponse{
			Address: server.Address,
			RoomID:  roomID,
		}
		// result, err := json.Marshal(&response)
		// if err != nil {
		// 	http.Error(w, "error marshaling response", http.StatusInternalServerError)
		// 	return
		// }
		log.Printf("New room created: %s at server: %s", response.RoomID, response.Address)
		w.WriteHeader(http.StatusOK)
		writeJSON(w, response)
	})
	mux.HandleFunc("/joinRoom", func(w http.ResponseWriter, r *http.Request) {
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
		response := PlayerResponse{
			Address: server,
			RoomID:  player.RoomID,
		}
		// result, err := json.Marshal(&response)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		log.Printf("New player joined room: %s at server: %s", response.RoomID, response.Address)
		w.WriteHeader(http.StatusOK)
		writeJSON(w, response)
	})
	mux.HandleFunc("/registerServer", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Errorf("Invalid Method").Error(), http.StatusMethodNotAllowed)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var server Server
		if err := decoder.Decode(&server); err != nil {
			http.Error(w, fmt.Errorf("Failed to decode body").Error(), http.StatusBadRequest)
			return
		}
		serverManager.AddServer(server)
		log.Printf("New server registered url: %s", server.Address)
		w.WriteHeader(http.StatusOK)
	})
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
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
		// http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}