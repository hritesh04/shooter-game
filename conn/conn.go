package conn

import (
	"github.com/hritesh04/shooter-game/types"
)

// type Connection struct {
// 	Device string
// 	// client pb.MovementEmitterClient
// 	Conn types.IConnections
// }

var clients = map[types.Device]func(string) types.IConnection{
	types.Desktop: NewDesktopGrpcClient,
	// types.Web:     NewJSGrpcClient,
}

func NewGrpcClient(address string, device types.Device) types.IConnection {
	return clients[device](address)
}
