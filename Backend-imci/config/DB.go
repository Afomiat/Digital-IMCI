package config



import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(env *Env) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), env.PostgresDSN)
	log.Println("Using DSN:", env.PostgresDSN)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Test the connection
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("Connected to Postgres successfully 🚀")
	return dbpool
}
