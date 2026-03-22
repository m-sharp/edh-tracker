package lib

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DBName = "pod_tracker"

	openErr = "error opening mysql connection: %w"

	maxConnTTL   = time.Minute * 2
	maxConnCount = 10
)

type quietLogger struct {
	logger.Interface
}

func (q quietLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err == gorm.ErrRecordNotFound {
		return
	}
	q.Interface.Trace(ctx, begin, fc, err)
}

type DBClient struct {
	log    *zap.Logger
	GormDb *gorm.DB
}

func NewDBClient(cfg *Config, log *zap.Logger) (*DBClient, error) {
	log = log.Named("DBClient")

	username, err := cfg.Get(DBUsername)
	if err != nil {
		return nil, fmt.Errorf(openErr, err)
	}
	password, err := cfg.Get(DBPass)
	if err != nil {
		return nil, fmt.Errorf(openErr, err)
	}
	host, err := cfg.Get(DBHost)
	if err != nil {
		return nil, fmt.Errorf(openErr, err)
	}
	port, err := cfg.Get(DBPort)
	if err != nil {
		return nil, fmt.Errorf(openErr, err)
	}

	log = log.With(
		zap.String("Username", username),
		zap.String("Host", host),
		zap.String("Port", port),
		zap.String("Database", DBName),
	)

	config := &mysql.Config{
		User:      username,
		Passwd:    password,
		Net:       "tcp",
		Addr:      fmt.Sprintf("%s:%v", host, port),
		DBName:    DBName,
		ParseTime: true,
	}

	log.Debug("Dialing mysql DB")
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf(openErr, err)
	}

	db.SetConnMaxLifetime(maxConnTTL)
	db.SetMaxOpenConns(maxConnCount)
	db.SetMaxIdleConns(maxConnCount)

	gormDB, err := gorm.Open(gormmysql.New(gormmysql.Config{Conn: db}), &gorm.Config{
		Logger: quietLogger{logger.Default.LogMode(logger.Warn)},
	})
	if err != nil {
		return nil, fmt.Errorf("error opening gorm connection: %w", err)
	}

	inst := &DBClient{log: log, GormDb: gormDB}

	if err := inst.CheckConnection(); err != nil {
		return nil, fmt.Errorf("DB connection check failed: %w", err)
	}

	return inst, nil
}

func (d *DBClient) CheckConnection() error {
	d.log.Debug("Pinging DB for health check...")
	sqlDB, err := d.GormDb.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

type DBError struct {
	inner error
	query string
}

func NewDBError(query string, innerErr error) *DBError {
	return &DBError{inner: innerErr, query: query}
}

func (d *DBError) Error() string {
	return fmt.Sprintf("failed to execute query %q: %s", d.query, d.inner)
}
