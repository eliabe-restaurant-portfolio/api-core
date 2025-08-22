package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/eliabe-portfolio/restaurant-app/internal/connections/configs"
	"github.com/eliabe-portfolio/restaurant-app/internal/envs"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connection struct {
	conn    *gorm.DB
	session *sql.DB
}

func (pc *Connection) Get() *gorm.DB {
	return pc.conn
}

func (pc *Connection) Close() error {
	if err := pc.session.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %v", err)
	}
	log.Println("✅ postgresql connection closed.")
	return nil
}

func Connect(config *configs.Config) (*Connection, error) {
	path := config.Postgres
	conn, err := gorm.Open(postgres.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(selectEnvModeByEnv()),
	})
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %v", err)
	}

	session, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("could not get database instance: %v", err)
	}

	session.SetMaxOpenConns(25)
	session.SetMaxIdleConns(25)
	session.SetConnMaxLifetime(5 * time.Minute)

	if err := session.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	log.Println("✅ postgresql connection established.")

	return &Connection{
		conn:    conn,
		session: session,
	}, nil
}

func selectEnvModeByEnv() logger.LogLevel {
	if envs.IsDev() {
		return logger.Info
	}
	return logger.Silent
}
