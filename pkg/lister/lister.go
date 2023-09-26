package lister

import (
	_sql "database/sql"
	"fmt"
	"strconv"
)

type Lister interface {
	FetchRows() ([]interface{}, error)
}

type DBLister struct {
	db      *_sql.DB
	sql     string
	idKey   string
	id      int
	resetId int
}

func NewDBLister(db *_sql.DB, sql, idKey string, id int) Lister {
	return &DBLister{
		db:      db,
		sql:     sql,
		idKey:   idKey,
		id:      id,
		resetId: id,
	}
}

func (l *DBLister) FetchRows() ([]interface{}, error) {
	sql := fmt.Sprintf(l.sql, l.id)
	rows, err := l.db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("lister.DBLister.FetchRows: %v (%s)", err, l.sql)
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("lister.DBLister.FetchRows: %v (%s)", err, l.sql)
	}

	var result []interface{}

	rowTmpVals := make([][]byte, len(cols))
	rowTmp := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		rowTmp[i] = &rowTmpVals[i]
	}

	for rows.Next() {
		if err := rows.Scan(rowTmp...); err != nil {
			return nil, fmt.Errorf("lister.DBLister.FetchRows: %v (%s)", err, l.sql)
		}
		row := make(map[string]string)
		for i := 0; i < len(cols); i++ {
			row[cols[i]] = string(rowTmpVals[i])
		}
		id, err := strconv.Atoi(row[l.idKey])
		if err != nil {
			return nil, fmt.Errorf("lister.DBLister.FetchRows: %v (%s)", err, l.sql)
		}
		l.id = id
		result = append(result, row)
	}

	if len(result) == 0 {
		l.id = l.resetId
	}
	return result, nil
}
