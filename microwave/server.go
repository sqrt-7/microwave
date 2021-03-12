package microwave

import (
	"net/http"

	"google.golang.org/grpc"
)

type ServerWrapper interface {
	Run(*Microwave)
}

type GRPCWrapper struct {
	Port   string
	Server *grpc.Server
}

type HTTPWrapper struct {
	Port   string
	Server *http.Server
}
