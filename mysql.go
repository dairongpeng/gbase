package gbase

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // register mysql driver
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"time"
)

var DB *bun.DB

func InitMysql() {
	// DSN:Data Source Name
	DSN := Cfg().GetString("mysql.dsn")
	db, err := newBunDB(DSN)
	if err != nil {
		panic(err)
	}
	DB = db
}

func newBunDB(dsn string) (*bun.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sql %s: %w", dsn, err)
	}
	bundb := bun.NewDB(db, mysqldialect.New())
	bundb.SetMaxIdleConns(Cfg().GetInt("mysql.maxIdleConn"))
	bundb.SetMaxOpenConns(Cfg().GetInt("mysql.maxOpenConn"))
	bundb.SetConnMaxIdleTime(time.Duration(Cfg().GetInt("mysql.maxIdleSec") * int(time.Second)))
	bundb.SetConnMaxLifetime(time.Duration(Cfg().GetInt("mysql.maxLifeSec") * int(time.Second)))
	if err := bundb.Ping(); err != nil {
		return nil, fmt.Errorf("ping db %s: %w", dsn, err)
	}
	bundb.AddQueryHook(&QueryHook{})
	return bundb, nil
}

type QueryHook struct{}

func (h *QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	if Cfg().GetString("log.level") == "DEBUG" || Cfg().GetString("log.level") == "debug" {
		Debug(ctx).Msg(event.Query)
	}
	return ctx
}

func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	cost := time.Since(event.StartTime) / time.Millisecond
	Debug(ctx).Int64("cost", int64(cost)).Msg(event.Query)
}
