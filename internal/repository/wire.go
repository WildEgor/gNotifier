package repository

import (
	mongo "github.com/WildEgor/gNotifier/internal/repository/mongo"
	"github.com/google/wire"
)

var RepositoriesSet = wire.NewSet(
	mongo.NewMongoClient,
	mongo.NewMongoDatabase,
	mongo.NewTokensRepository,
	wire.Bind(new(mongo.ITokensRepository), new(*mongo.TokensRepository)),
)
