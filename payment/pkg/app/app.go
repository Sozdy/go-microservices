package app

import (
	"google.golang.org/grpc"
)

func Interceptors() []grpc.ServerOption {
	return nil
}

func RegisterServices(server *grpc.Server) {

}
