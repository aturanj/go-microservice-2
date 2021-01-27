package main

import (
	"fmt"
	"go-url-shortener/shortener"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	api "go-url-shortener/api"
	mng "go-url-shortener/repository/mongodb"
	rds "go-url-shortener/repository/redis"
)

func main() {

	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := api.NewHandler(service)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/{code}", handler.Get)
	router.Post("/", handler.Post)

	errors := make(chan error, 2)

	go func() {
		fmt.Println("Listening port : 8000")
		errors <- http.ListenAndServe(httpPort(), router)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errors <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errors)
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
