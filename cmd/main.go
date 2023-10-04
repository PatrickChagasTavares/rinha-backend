package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/patrickchagastavares/rinha-backend/internal/controllers"
	"github.com/patrickchagastavares/rinha-backend/internal/handlers"
	"github.com/patrickchagastavares/rinha-backend/internal/repositories"
	"github.com/patrickchagastavares/rinha-backend/internal/services"
	"github.com/patrickchagastavares/rinha-backend/pkg/httpRouter"
	"github.com/patrickchagastavares/rinha-backend/pkg/logger"
	"github.com/patrickchagastavares/rinha-backend/pkg/migration"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	godotenv.Load(".env")

	migration.RunMigrations(os.Getenv("DATABASE_WRITE_URL"))
	fmt.Println("os.Getenv(DATABASE_WRITE_URL) -> ", os.Getenv("DATABASE_WRITE_URL"))

	var (
		log          = logger.NewLogrusLogger()
		repositories = repositories.New(repositories.Options{
			WriterSqlx: sqlx.MustConnect("postgres", os.Getenv("DATABASE_WRITE_URL")),
			ReaderSqlx: sqlx.MustConnect("postgres", os.Getenv("DATABASE_READ_URL")),
			Log:        log,
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

	log.Info("start serve in port:", os.Getenv("PORT"))
	if err := router.Server(os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}

}
