package adapters

import "github.com/google/wire"

var AdaptersSet = wire.NewSet(
	NewHealthCheckAdapter,
)
