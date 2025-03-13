package postgres

import (
	"context"
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

// connection string exmaple:
// postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
// fmt.Sprintf(
//
//	"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
//	config.Username,
//	config.Password,
//	config.Host,
//	config.Port,
//	config.Database,
//	config.SSLMode,
//
// )
func NewPostgresConnection(config config.PostgresConfig) *PostgresConnection {

	poolConfig, err := pgxpool.ParseConfig(config.PostgresUrl)
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

	dbPool, err := pgxpool.NewWithConfig(ctx, &pgxpool.Config{})
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

func (r *PostgresConnection) Close() {
	r.dbPool.Close()
}
