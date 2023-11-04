package mongo

import (
	"context"
	"github.com/pkg/errors"
	"strings"
	"time"

	"github.com/WildEgor/gNotifier/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName    string = "notification_tokens"
	mongoQueryTimeout        = 10 * time.Second
)

type TokensFilter struct {
	SubId    string
	Token    string
	Platform string
}

type ITokensRepository interface {
	FindSub(f *TokensFilter) (*models.SubTokenModel, error)
	UpsertToken(m *models.SubTokenCreateModel) (*models.SubTokenModel, error)
}

type TokensRepository struct {
	collection *mongo.Collection
}

func NewTokensRepository(db *mongo.Database) (*TokensRepository, error) {
	return &TokensRepository{
		collection: db.Collection(collectionName),
	}, nil
}

func (r *TokensRepository) FindSub(f *TokensFilter) (*models.SubTokenModel, error) {
	var result *models.SubTokenModel
	var (
		andQuery []bson.M
		query    []bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	if len(f.Platform) > 0 {
		query = append(query, bson.M{"platform": bson.M{"$eq": f.Platform}})
	}

	if len(f.SubId) > 0 {
		query = append(query, bson.M{"sub_id": bson.M{"$eq": f.SubId}})
	}

	if len(f.Token) > 0 {
		query = append(query, bson.M{"tokens": bson.M{"$in": []string{f.Token}}})
	}

	andQuery = append(andQuery, bson.M{"$and": func() []bson.M {
		if query != nil {
			if len(query) > 0 {
				return query
			}
		}
		return []bson.M{}
	}()})

	err := r.collection.FindOne(ctx, andQuery)
	if err != nil {
		if errors.Is(err.Err(), mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err.Err()
	}

	dErr := err.Decode(&result)
	if dErr != nil {
		return nil, dErr
	}

	return result, nil
}

func (r *TokensRepository) UpsertToken(m *models.SubTokenCreateModel) (*models.SubTokenModel, error) {
	var existedResult *models.SubTokenModel
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	err := r.collection.FindOne(ctx, bson.M{
		"$and": bson.M{
			"sub_id": bson.M{
				"$eq": m.SubID,
			},
		},
	})
	if err != nil {
		if errors.Is(err.Err(), mongo.ErrNoDocuments) {
		}
	} else {
		err := err.Decode(&existedResult)
		if err != nil {
			return nil, err
		}
	}

	if existedResult != nil {
		findQuery := bson.M{
			"_id":          existedResult.ID,
			"tokens.token": m.Token.Token,
		}
		var updateQuery bson.M
		for _, v := range existedResult.Tokens {
			if strings.EqualFold(v.Token, m.Token.Token) {
				updateQuery = bson.M{
					"$set": bson.M{
						"tokens.$.updated_at": time.Now(),
					},
				}
			}
		}

		if updateQuery == nil {
			updateQuery = bson.M{
				"$push": bson.M{
					"tokens": &models.TokenModel{
						Platform:  m.Token.Platform,
						Token:     m.Token.Token,
						UpdatedAt: time.Now(),
						CreatedAt: time.Now(),
					},
				},
			}
		}

		_, err := r.collection.UpdateOne(ctx, findQuery, updateQuery)
		if err != nil {
			return nil, err
		}
	}

	newTokenModel := &models.TokenModel{
		Platform:  m.Token.Platform,
		Token:     m.Token.Token,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	newSubModel := &models.SubTokenModel{
		SubID:  m.SubID,
		Tokens: []*models.TokenModel{newTokenModel},
	}

	saveResult, er := r.collection.InsertOne(context.TODO(), &newSubModel)
	if er != nil {
		return nil, er
	}

	newSubModel.ID = saveResult.InsertedID.(primitive.ObjectID)
	return newSubModel, nil
}
