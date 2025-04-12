package mysql_client

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"log"
	"time"
)

type (
	MySQL interface {
		Close()
		Ping() error
		Client() *sql.DB
		Migration(ctx context.Context) error
	}

	mysqlClient struct {
		db *sql.DB
	}
)

func (m *mysqlClient) Close() {
	m.db.Close()
}

func (m *mysqlClient) Ping() error {
	return m.db.Ping()
}

func (m *mysqlClient) Client() *sql.DB {
	return m.db
}

var _ MySQL = &mysqlClient{}

func NewMySQLClient(config *configs.Config) (*mysqlClient, error) {
	mysql := config.Config.MySQL
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysql.Username, mysql.Password, mysql.Host, mysql.Port, mysql.Database)
	log.Println(url)
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	connPool := mysql.ConnPool
	db.SetMaxIdleConns(connPool.MaxIdle)
	db.SetMaxOpenConns(connPool.ConnLimit)
	db.SetConnMaxIdleTime(connPool.IdleTime * time.Second)

	return &mysqlClient{
		db: db,
	}, nil
}
