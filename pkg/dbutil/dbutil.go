package dbutil

import (
	"database/sql"
	"sync"
	"time"
)

func FetchAll(rows *sql.Rows, err error) ([]map[string]string, error) {
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var res []map[string]string

	rowTmpVals := make([][]byte, len(cols))
	rowTmp := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		rowTmp[i] = &rowTmpVals[i]
	}

	for rows.Next() {
		if err := rows.Scan(rowTmp...); err != nil {
			return nil, err
		}
		row := make(map[string]string)
		for i := 0; i < len(cols); i++ {
			row[cols[i]] = string(rowTmpVals[i])
		}
		res = append(res, row)
	}

	return res, nil
}

func LastInsertId(result sql.Result, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	ret, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func RowsAffected(result sql.Result, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	ret, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return ret, nil
}

type DBWrapper struct {
	driverName      string
	dataSourceName  string
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime int
	db              *sql.DB
	mu              sync.RWMutex
}

func NewDBRaw(db *sql.DB) *DBWrapper {
	return &DBWrapper{
		db: db,
	}
}

func NewDB(driverName, dataSourceName string, maxOpenConns, maxIdleConns, connMaxLifetime int) *DBWrapper {
	return &DBWrapper{
		driverName:      driverName,
		dataSourceName:  dataSourceName,
		maxOpenConns:    maxOpenConns,
		maxIdleConns:    maxIdleConns,
		connMaxLifetime: connMaxLifetime,
	}
}

func (dbw *DBWrapper) Open() (*sql.DB, error) {
	dbw.mu.RLock()
	if dbw.db != nil {
		dbw.mu.RUnlock()
		return dbw.db, nil
	}
	dbw.mu.RUnlock()
	dbw.mu.Lock()
	defer dbw.mu.Unlock()
	db, err := sql.Open(dbw.driverName, dbw.dataSourceName)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(dbw.maxOpenConns)
	db.SetMaxIdleConns(dbw.maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(dbw.connMaxLifetime) * time.Second)
	dbw.db = db
	return db, nil
}

func (dbw *DBWrapper) Prepare(query string) (*sql.Stmt, error) {
	db, err := dbw.Open()
	if err != nil {
		return nil, err
	}
	return db.Prepare(query)
}

func (dbw *DBWrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := dbw.Open()
	if err != nil {
		return nil, err
	}
	return db.Query(query, args...)
}

func (dbw *DBWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	db, err := dbw.Open()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, args...)
}

func (dbw *DBWrapper) Begin() (*sql.Tx, error) {
	db, err := dbw.Open()
	if err != nil {
		return nil, err
	}
	return db.Begin()
}
