package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/IbnBaqqi/gitchat/internal/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// DB wraps the database connection pool & queries
type DB struct {
	*sql.DB
	*Queries
	log *zap.Logger
}

// Tx wraps a database transaction
type Tx struct {
	*sql.Tx
	log *zap.Logger
}

// Connect establishes a connection to the database
func Connect(ctx context.Context, cfg *config.DBConfig, log *zap.Logger) (*DB, error) {
	
	driver := "postgres"
	
	dbConn, err := sql.Open(driver, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := dbConn.PingContext(ctx); err != nil {
		if err = dbConn.Close(); err != nil {
			log.Error("failed to close database connection", zap.Error(err))
		}
		log.Error("failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	dbConn.SetMaxOpenConns(cfg.MaxOpenConns) // 25
	dbConn.SetMaxIdleConns(cfg.MaxIdleConns) // 5
	dbConn.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 5minutes

	queries := New(dbConn)

	return &DB{
		DB:      dbConn,
		Queries: queries,
		log: log,
	}, nil
}

// Close closes the database connection and logs the closure.
func (db *DB) Close() error {
	db.log.Info("closing database connection")
	return db.DB.Close()
}

// BeginTx starts a new database transaction with the specified isolation level
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		db.log.Error("failed to begin transaction", zap.Error(err))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	db.log.Debug("transaction started")
	return &Tx{
		Tx: tx,
		log: db.log,
	}, nil
}

// Commit commits the transaction
func (tx *Tx) Commit() error {
	if err := tx.Tx.Commit(); err != nil {
		tx.log.Error("failed to commit transaction", zap.Error(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	tx.log.Debug("transaction committed")
	return nil
}

// Rollback rolls back the transaction
func (tx *Tx) Rollback() error {
	if err := tx.Tx.Rollback(); err != nil {
		if errors.Is(err, sql.ErrTxDone) {
			tx.log.Debug("transaction already closed, ignoring rollback")
			return nil
		}
		tx.log.Error("failed to rollback transaction", zap.Error(err))
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	tx.log.Debug("transaction rolled back")
	return nil
}