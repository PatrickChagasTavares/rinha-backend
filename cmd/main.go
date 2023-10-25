package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/patrickchagastavares/rinha-backend/internal/controllers"
	"github.com/patrickchagastavares/rinha-backend/internal/handlers"
	"github.com/patrickchagastavares/rinha-backend/internal/repositories"
	"github.com/patrickchagastavares/rinha-backend/internal/services"
	"github.com/patrickchagastavares/rinha-backend/pkg/httpRouter"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
)

func main() {
	godotenv.Load(".env")

	var (
		log          = logger.NewLogrusLogger()
		repositories = repositories.New(repositories.Options{
			DB_URL: os.Getenv("DATABASE_URL"),
			Log:    log,
		})
		services = services.New(services.Options{
			Repo: repositories,
			Log:  log,
		})
		controllers = controllers.New(controllers.Options{
			Srv: services,
			Log: log,
		})
		router = httpRouter.NewGinRouter()
	)

	handlers.NewRouter(handlers.Options{
		Router: router,
		Ctrl:   controllers,
	})

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8000"
	}
	log.Info("start serve in port:", port)
	if err := router.Server(port); err != nil {
		log.Fatal(err)
	}

}
