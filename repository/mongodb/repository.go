package mongodb

import (
	"context"
	"go-url-shortener/handler"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	mongodbClient, err := newMongodbClient(mongodbURL, mongodbTimeout)

	if err != nil {
		return nil, errors.Wrap(err, "repository.NewMongodbRepository")
	}

	mongodbRepository := &mongodbRepository{
		timeout:  time.Duration(mongodbTimeout) * time.Second,
		database: mongodb,
	}
	mongodbRepository.client = mongodbClient

	return mongodbRepository, nil
}

func (r *mongodbRepository) Find(code string) (*handler.Redirect, error) {

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("redirects")

	filter := bson.M{"code": code}

	redirect := &handler.Redirect{}

	err := collection.FindOne(ctx, filter).Decode(&redirect)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrap(handler.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	return redirect, nil
}

func (r *mongodbRepository) Store(redirect *handler.Redirect) error {

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("redirects")

	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":      redirect.Code,
			"url":       redirect.URL,
			"credit_at": redirect.CreatedAt,
		},
	)

	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}

	return nil
}
