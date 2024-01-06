//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
)

var ServerSet = wire.NewSet(AppSet)

func NewServer() (*Server, error) {
	wire.Build(ServerSet)
	return nil, nil
}
