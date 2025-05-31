package app

import (
	"context"
	"os"

	"github.com/Sanchir01/wallets/internal/config"
	"github.com/Sanchir01/wallets/pkg/lib/db/connect"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	PrimaryDB *pgxpool.Pool
	RedisDB   *redis.Client
}

func NewDataBases(cfg *config.Config) (*Database, error) {
	pgxdb, err := connect.PGXNew(context.Background(), cfg.DB.User, cfg.DB.Host, cfg.DB.Database, cfg.DB.Port, cfg.DB.MaxAttempts)
	if err != nil {
		return nil, err
	}
	redisdb, err := connect.RedisConnect(context.Background(), cfg.RedisDB.Host, cfg.RedisDB.Port,
		os.Getenv("REDIS_PASSWORD"), cfg.Env,
		cfg.RedisDB.DBNumber, cfg.RedisDB.Retries)
	if err != nil {
		return nil, err
	}
	return &Database{PrimaryDB: pgxdb, RedisDB: redisdb}, nil
}

func (databases *Database) Close() error {
	databases.PrimaryDB.Close()

	return nil
}
