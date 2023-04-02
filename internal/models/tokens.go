package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenModel struct {
	Platform  string    `json:"platform"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubTokenModel struct {
	ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	SubID  string             `json:"sub_id"`
	Tokens []*TokenModel      `json:"tokens"`
}

type SubTokenCreateModel struct {
	SubID string      `json:"sub_id"`
	Token *TokenModel `json:"token"`
}
