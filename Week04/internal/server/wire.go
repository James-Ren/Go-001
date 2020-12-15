// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package server

import (
	"github.com/James-Ren/Go-001/tree/main/Week04/internal/service"
	"github.com/google/wire"
)

func InitializeServer() (*Server, func(), error) {
	wire.Build(NewServer, service.Provider)
	return nil, nil, nil
}
