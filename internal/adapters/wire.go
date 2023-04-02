package adapters

import (
	"github.com/google/wire"
)

var AdaptersSet = wire.NewSet(
	NewHealthCheckAdapter,
	NewFCMAdapter,
	NewAPNAdapter,
	NewSMSAdapter,
	NewSMSAdapter,
)
