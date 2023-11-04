package mongo

import (
	"github.com/WildEgor/gNotifier/internal/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(
	cfg *configs.MongoConfig,
) (*mongo.Client, error) {
	opts := options.Client()
	opts.Hosts = append(opts.Hosts, cfg.GetHost())
	opts.Auth = &options.Credential{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	return mongo.NewClient(opts)
}

func NewMongoDatabase(
	client *mongo.Client,
) *mongo.Database {
	return client.Database("test")
}
