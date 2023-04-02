package adapters

import (
	"github.com/google/wire"
)

var AdaptersSet = wire.NewSet(
	NewHealthCheckAdapter,
	NewFCMAdapter,
	wire.Bind(new(IFCMAdapter), new(*FCMAdapter)),
	NewAPNAdapter,
	wire.Bind(new(IAPNAdapter), new(*APNAdapter)),
	NewSMSAdapter,
	wire.Bind(new(ISMSAdapter), new(*SMSAdapter)),
	NewSMTPAdapter,
	wire.Bind(new(ISMTPAdapter), new(*SMTPAdapter)),
)
