package postgresql_sqlboiler

import (
	"database/sql"
	"fmt"
	"github.com/hdget/common/types"
	_ "github.com/jackc/pgx/v5/stdlib" // 替代原来的pq驱动
	"time"
)

type psqlClient struct {
	*sql.DB
}

const (
	// DSN (Data Type NickName): username:password@address/dbname?param=value
	dsnTemplate = "postgres://%s:%s@%s:%d/%s?TimeZone=Asia/Shanghai"
)

func newClient(c *psqlConfig) (types.DbClient, error) {
	// 构造连接参数
	dsn := fmt.Sprintf(dsnTemplate, c.User, c.Password, c.Host, c.Port, c.Database)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(c.MaxOpenConn)                     // 最大活跃连接数（对应PgBouncer的max_client_conn）
	db.SetMaxIdleConns(int(0.2 * float32(c.MaxOpenConn))) // 最大空闲连接数（建议值=0.2*MaxOpenConns)
	if !c.UsePgBouncer {
		db.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Second) // 应小于PostgreSQL的idle_in_transaction_session_timeout, 默认为60m
	}

	return &psqlClient{DB: db}, nil
}

func (m psqlClient) Close() error {
	return m.DB.Close()
}

func (m psqlClient) Get(dest interface{}, query string, args ...interface{}) error {
	return nil
}

func (m psqlClient) Select(dest interface{}, query string, args ...interface{}) error {
	return nil
}

func (m psqlClient) Rebind(query string) string {
	return ""
}
