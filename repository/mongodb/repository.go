package mongodb

import (
	"context"
	"go-url-shortener/handler"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/pkg/errors"
)

type mongodbRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongodbClient(mongodbURL string, mongodbTimeout int) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongodbTimeout)*time.Second)
	defer cancel()

	mongodbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbURL))

	if err != nil {
		return nil, err
	}

	err = mongodbClient.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, err
	}

	return mongodbClient, nil
}

//NewMongodbRepository returns RedirectionRepository
func NewMongodbRepository(mongodbURL string, mongodb string, mongodbTimeout int) (handler.RedirectRepository, error) {

	mongodbRepository := &mongodbRepository{
		timeout:  time.Duration(mongodbTimeout) * time.Second,
		database: mongodb,
	}

	mongodbClient, err := newMongodbClient(mongodbURL, mongodbTimeout)

	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongodbRepository")
	}

	mongodbRepository.client = mongodbClient

	return mongodbRepository, nil
}

func (r *mongodbRepository) Find(code string) (*handler.Redirect, error) {

}

func (r *mongodbRepository) Store(redirect *handler.Redirect) err {

}
