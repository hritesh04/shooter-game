package conn

import (
	"github.com/hritesh04/shooter-game/conn/grpc"
	"github.com/hritesh04/shooter-game/conn/rest"
	"github.com/hritesh04/shooter-game/types"
)

var grpcClient = map[types.Device]func(string) types.IConnection{
	types.Desktop: grpc.NewGrpcDesktopClient,
}

var restClient = map[types.Device]func(string) types.IConnection{
	types.Desktop: rest.NewRestDesktopClient,
}

func NewGrpcClient(address string, device types.Device) types.IConnection {
	return grpcClient[device](address)
}

func NewRestClient(address string, device types.Device) types.IConnection {
	return restClient[device](address)
}
