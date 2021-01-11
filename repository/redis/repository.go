package redis

import (
	"fmt"
	"go-url-shortener/handler"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	_, err = client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}

//NewRedisRepository returns RedirectRepository
func NewRedisRepository(redisURL string) (handler.RedirectRepository, error) {

	client, err := newRedisClient(redisURL)

	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}

	repo := &redisRepository{}
	repo.client = client

	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*handler.Redirect, error) {

	key := r.generateKey(code)

	data, err := r.client.HGetAll(key).Result()

	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	if len(data) == 0 {
		return nil, errors.Wrap(handler.ErrRedirectNotFound, "repository.Redirect.Find")
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)

	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	redirect := &handler.Redirect{}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt

	return redirect, nil
}

func (r *redisRepository) Store(redirect *handler.Redirect) error {

	key := r.generateKey(redirect.Code)

	data := map[string]interface{}{
		"code":       redirect.Code,
		"URL":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	_, err := r.client.HMSet(key, data).Result()

	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}

	return nil
}
