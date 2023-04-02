package config

import "github.com/google/wire"

var ConfigSet = wire.NewSet(
	NewAppConfig,
	NewAMQPConfig,
	NewFCMConfig,
	NewAPNConfig,
	NewSMSConfig,
	NewSMTPConfig,
)
