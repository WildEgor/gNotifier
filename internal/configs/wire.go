package configs

import "github.com/google/wire"

var ConfigSet = wire.NewSet(
	NewConfigurator,
	NewAppConfig,
	NewAMQPConfig,
	NewFCMConfig,
	NewAPNConfig,
	NewSMSConfig,
	NewSMTPConfig,
	NewMongoConfig,
)
