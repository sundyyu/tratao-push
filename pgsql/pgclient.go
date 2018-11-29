package pgclient

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"sync"
	"xcurrency-push/config"
	"xcurrency-push/util"
)

var sqlDB *sql.DB
var dbLock *sync.Mutex = new(sync.Mutex)
var connChan chan *sql.DB

func initConn() error {
	cfg := config.GetConfig()
	host := cfg.GetString("postgresql.ip")
	port := cfg.GetInt("postgresql.port")
	dbname := cfg.GetString("postgresql.dbname")
	user := cfg.GetString("postgresql.user")
	password := cfg.GetString("postgresql.password")

	pqsql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	connChan = make(chan *sql.DB, 10)
	for i := 0; i < 10; i++ {
		if db, err := sql.Open("postgres", pqsql); err == nil {
			connChan <- db
		} else {
			return err
		}
	}
	util.LogInfo("PostgreSQL conn init.")
	return nil
}

func GetConn() (*sql.DB, error) {
	if connChan == nil {
		dbLock.Lock()
		defer dbLock.Unlock()
		if connChan == nil {
			if err := initConn(); err != nil {
				util.LogError(err)
				return nil, err
			}
		}
	}
	db := <-connChan
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func ReleaseConn(db *sql.DB) {
	connChan <- db
}

func CloseConn() {
	for i := 0; i < 10; i++ {
		db := <-connChan
		db.Close()
	}
}

//插入、新增、删除数据
func ExecBySQL(pgsql string, data ...interface{}) error {
	db, err := GetConn()
	if err != nil {
		return err
	}
	defer ReleaseConn(db)
	stmt, err := db.Prepare(pgsql)
	if err != nil {
		util.LogInfo(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(data...)
	if err != nil {
		util.LogInfo(err)
		return err
	}
	return nil
}
