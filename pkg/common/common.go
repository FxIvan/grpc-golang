package common

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectToDatabase(ctx context.Context, dbConnectionString string) (*pgxpool.Pool, error) {
	var dbPool *pgxpool.Pool
	var err error
	//Hago una variable que acumula el reintento de conexion
	retryCount := 0

	for retryCount < 5 {
		dbPool, err = pgxpool.Connect(ctx, dbConnectionString)
		if err == nil {
			break
		}
		log.Printf("Failed to Connect database, Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		retryCount++
	}

	if err != nil {
		log.Printf("Ran out of retries to connect to database (5)")
	}

	log.Printf("Connect to the database.")
	return dbPool, nil

}
