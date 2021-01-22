package main

import (
	"fmt"
	"go-url-shortener/shortener"
	"log"
	"os"
	"strconv"

	api "go-url-shortener/api"
	mng "go-url-shortener/repository/mongodb"
	rds "go-url-shortener/repository/redis"
)

func main() {

	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := api.NewHandler(service)

}

func httpPort() string {

	port := "8000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	return fmt.Sprintf(":%s", port)
}

func chooseRepo() shortener.RedirectRepository {

	switch os.Getenv("URL-DB") {
	case "redis":

		redisURL := os.Getenv("REDIS_URL")

		repo, err := rds.NewRedisRepository(redisURL)

		if err != nil {
			log.Fatal(err)
		}
		return repo

	case "mongo":

		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))

		repo, err := mng.NewMongodbRepository(mongoURL, mongodb, mongoTimeout)

		if err != nil {
			log.Fatal(err)
		}

		return repo
	}

	return nil
}
