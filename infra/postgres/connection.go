package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nsaltun/todolist-service/config"
	"go.uber.org/zap"
)

const applicationName = "todolist-service"

type PostgresConnection struct {
	config config.PostgresConfig
	dbPool *pgxpool.Pool
}

func NewPostgresConnection(config config.PostgresConfig) *PostgresConnection {
	poolConfig, err := pgxpool.ParseConfig(connectionString(config))
	if err != nil {
		zap.L().Fatal("Unable to parse pool config", zap.Error(err))
		return nil
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 10 * time.Minute
	poolConfig.HealthCheckPeriod = 30 * time.Second
	poolConfig.MaxConnLifetimeJitter = 2 * time.Minute
	poolConfig.ConnConfig.ConnectTimeout = 5 * time.Second // Connection timeout
	poolConfig.ConnConfig.RuntimeParams = map[string]string{
		"application_name":  applicationName, // Identify your application in pg_stat_activity
		"statement_timeout": "30000",         // Statement timeout in milliseconds (30s)
		"search_path":       "public",        // Default schema
	}

	// Create the connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		zap.L().Fatal("Unable to create connection pool", zap.Error(err))
		return nil
	}

	//Verify the connection
	if err := dbPool.Ping(ctx); err != nil {
		zap.L().Fatal("Unable to ping database", zap.Error(err))
		dbPool.Close()
		return nil
	}

	zap.L().Info("Successfully connected to database",
		zap.String("host", poolConfig.ConnConfig.Host),
		zap.String("database", poolConfig.ConnConfig.Database),
		zap.Int32("max_connections", poolConfig.MaxConns),
	)

	return &PostgresConnection{config: config, dbPool: dbPool}
}

func connectionString(c config.PostgresConfig) string {
	//postgres://user:password@host:port/database?sslmode=disable
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		"disable",
	)
}

func (r *PostgresConnection) Close() {
	r.dbPool.Close()
}
